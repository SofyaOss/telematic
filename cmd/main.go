package main

import (
	"context"
	"errors"
	"flag"
	"github.com/jackc/pgx/v4"
	"net"
	"practice/internal/grpc_server"

	"strconv"

	//"encoding/json"
	//"fmt"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"practice/internal/generator"
	myKafka "practice/internal/kafka"
	"practice/internal/redis"
	"practice/storage"
	"practice/storage/postgres"
	"time"
)

func main() {
	log.Println("Start telematic service...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := postgres.New(ctx, os.Getenv("DATABASE_URL")) // подключение к бд
	if err != nil {
		log.Fatalf("Could not connect to postgres db: %s", err)
	}

	err = db.DropTable() // удаление старой таблицы
	if err != nil {
		log.Println("Could not drop the table:", err)
	}

	err = db.CreateTable() // создание новой таблицы
	if err != nil {
		log.Println("Could not create the table", err)
	}

	// determine the port the container is listening on
	addr := "redis:6379"

	newRedis := redis.New(addr)                    // создание клиента редис
	newKafkaProducer, err := myKafka.NewProducer() // создание продюсера кафки
	if err != nil {
		log.Fatalf("Could not connect to kafka: %s", err)
	}

	kafkaCh := make(chan *storage.Car)
	amount, err := strconv.Atoi(os.Getenv("TRANSPORT_AMOUNT")) // количество машин
	if err != nil {
		log.Fatal("Transport amount must be an integer")
	}
	for i := 0; i < amount; i++ {
		go generator.Generate(i, kafkaCh)
	}

	key := 0
	go func() {
		for {
			val, ok := <-kafkaCh
			if ok == false {
				close(kafkaCh)
				break // exit break loop
			} else {
				val.ID = key
				err = newRedis.AddToRedis(val, key) // добавление в редис
				if err != nil {
					log.Fatalf("Could not add element to redis: %s", err)
				}
				err = db.AddData(val) // добавиление в бд
				if err != nil && !errors.Is(pgx.ErrNoRows, err) {
					log.Fatalf("Could not add element to db: %s", err)
				}
				err = myKafka.Produce(newKafkaProducer, val) // отправка сообщения в кафку
				if err != nil {
					log.Fatalf("Could not add element to kafka: %s", err)
				}
				key++
			}
		}
	}()

	// создание консюмера кафки
	/*
		log.Println("start consumer")
		time.Sleep(3 * time.Second)

		config2 := &kafka.ConfigMap{
			"bootstrap.servers": "kafka:9092",
			"group.id":          "myGroup",
			"auto.offset.reset": "earliest",
		}
		consumer, err := kafka.NewConsumer(config2)
		if err != nil {
			panic(err)
		}
		consumer.SubscribeTopics([]string{"telematicTopic"}, nil)
		go func() {
			for {
				msg, err := consumer.ReadMessage(-1)
				if err == nil {
					log.Printf("Recieved message: %s\n", string(msg.Value))
				} else {
					log.Printf("Error while consuming message %v(%v)\n", err, msg)
				}
			}
			consumer.Close()
		}()
	*/

	// создание grpc сервера
	var gRPCAddr string
	flag.StringVar(&gRPCAddr, "grpc-addr", "localhost:8000", "Set the grpc address")

	lis, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Failed to listen to port: %s", err)
	}
	srv := grpc_server.New(db)
	if err := srv.Grpc.Serve(lis); err != nil {
		log.Fatalf("Could not start grpc server: %s", err)
	}
}
