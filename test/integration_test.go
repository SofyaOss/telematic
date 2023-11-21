package test

import (
	"context"
	"flag"
	//"github.com/go-redis/redis"
	"net"
	"practice/internal/grpc_server"
	"practice/internal/kafka"
	"practice/internal/redis"
	//"strconv"

	//"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"practice/internal/generator"
	//myRedis "practice/internal/redis"
	"practice/storage"
	"practice/storage/postgres"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	redisRes, err := pool.Run("redis", "5-alpine", nil)
	if err != nil {
		log.Fatalf("Failed to start redis: %+v", err)
	}
	defer pool.Purge(redisRes)

	// determine the port the container is listening on
	addr := net.JoinHostPort("localhost", redisRes.GetPort("6379/tcp"))

	optsDB := dockertest.RunOptions{
		Name:         "postgres-test",
		Repository:   "postgres",
		Tag:          "latest",
		Env:          []string{"POSTGRES_PASSWORD=postgres", "POSTGRES_USER=postgres", "POSTGRES_DB=testdb"},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5432"},
			},
		},
	}

	resourceDB, err := pool.RunWithOptions(&optsDB)
	if err != nil {
		log.Fatalf("could not start resourse: %s", err)
	}
	defer pool.Purge(resourceDB)

	zookeeperRes, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "zookeeper-test",
		Repository:   "confluentinc/cp-zookeeper",
		Tag:          "latest",
		Hostname:     "zookeeper",
		ExposedPorts: []string{"2182"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"2182": {
				{HostIP: "0.0.0.0", HostPort: "2182"},
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not start zookeeper: %s", err)
	}
	defer pool.Purge(zookeeperRes)

	kafkaRes, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "kafka-test",
		Repository: "confluentinc/cp-kafka",
		Tag:        "latest",
		Hostname:   "kafka",
		Env: []string{
			"KAFKA_BROKER_ID: 1",
			"KAFKA_ZOOKEEPER_CONNECT: zookeeper:2182",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1",
			"KAFKA_ADVERTISED_LISTENERS: INSIDE://:9092,OUTSIDE://:9093",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9093/tcp": {{HostIP: "localhost", HostPort: "9093/tcp"}},
		},
		ExposedPorts: []string{"9093/tcp"},
	})
	if err != nil {
		log.Fatalf("Could not start kafka: %s", err)
	}
	defer pool.Purge(kafkaRes)

	var db *postgres.TelematicDB // wait for the container to be ready
	err = pool.Retry(func() error {
		client := redis.New(addr)
		//client := myRedis.New()
		//defer client.Close()

		newKafkaProducer, err := kafka.NewProducer() // создание продюсера кафки
		if err != nil {
			log.Fatalf("Could not connect to kafka: %s", err)
		}

		dbConnStr := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=testdb password=postgres sslmode=disable")
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		db, err = postgres.New(ctx, dbConnStr) // подключение к базе данных
		if err != nil {
			log.Println("db not ready yet")
			return err
		}
		err = db.DropTable()
		if err != nil {
			log.Fatalf("Could not drop table: %s", err)
		}
		err = db.CreateTable()
		if err != nil {
			log.Fatalf("Could not create table: %s", err)
		}

		kafkaCh := make(chan *storage.Car)
		go generator.Generate(1, kafkaCh)
		key := 0
		go func() {
			for {
				val, ok := <-kafkaCh
				if ok == false {
					close(kafkaCh)
					break // exit break loop
				} else {
					err = db.AddData(val)
					if err != nil && !errors.Is(pgx.ErrNoRows, err) {
						log.Fatalf("Could not add element to db: %s", err)
					}
					//mes, err := json.Marshal(val)
					//if err != nil {
					//	log.Fatalf("Could not convert to json: %s", err)
					//}
					////err = myRedis.AddToRedis(client, val, key)
					//if err != nil {
					//	log.Fatalf("could not add element to redis: %s", err)
					//}
					err = client.AddToRedis(val, key)
					if err != nil {
						log.Fatalf("Could not add to redis: %s", err)
					}
					err = kafka.Produce(newKafkaProducer, val)
					if err != nil {
						log.Fatalf("Could not send message to kafka: %s", err)
					}
					key++
				}
			}
		}()

		var gRPCAddr string
		flag.StringVar(&gRPCAddr, "grpc-addr", "localhost:8000", "Set the grpc address")

		lis, err := net.Listen("tcp", ":8000")
		if err != nil {
			log.Fatalf("Could not listen grpc port: %s", err)
		}
		srv := grpc_server.New(db)
		if err := srv.Grpc.Serve(lis); err != nil {
			log.Fatalf("Could not start grpc server: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to ping Redis: %+v", err)
	}

	code := m.Run()
	if err = pool.Purge(redisRes); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err = pool.Purge(kafkaRes); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err = pool.Purge(zookeeperRes); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err = pool.Purge(resourceDB); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}
