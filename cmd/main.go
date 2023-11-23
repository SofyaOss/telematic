package main

import (
	"context"
	"log"
	"os"
	"time"

	"practice/internal/grpc_server"
	myKafka "practice/internal/kafka"
	"practice/internal/redis"
	"practice/storage"
	"practice/storage/postgres"
)

type Application struct {
	Channel  chan *storage.Car
	Redis    *redis.Client
	Postgres *postgres.TelematicDB
	Kafka    *myKafka.Producer
	GRPC     *grpc_server.Server
}

func main() {
	log.Println("Start telematic service...")

	conf := getConfig() // настройки сервиса

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := postgres.New(ctx, os.Getenv("DATABASE_URL")) // подключение к бд
	if err != nil {
		log.Fatalf("Could not connect to postgres db: %s", err)
	}

	newRedis := redis.New(getRedisConf(conf)) // создание клиента редис

	newKafkaProducer, err := myKafka.NewProducer(getKafkaConf(conf)) // создание продюсера кафки
	if err != nil {
		log.Fatalf("Could not connect to kafka: %s", err)
	}

	telematicCh := make(chan *storage.Car) // канал для передачи телематики

	app := &Application{
		Channel:  telematicCh,
		Redis:    newRedis,
		Postgres: db,
		Kafka:    newKafkaProducer,
	}

	app.refreshTelematicTable(ctx)
	app.startApp(ctx)
	app.startServer(getGRPCConf(conf))

	//lis, err := net.Listen("tcp", getGRPCConf(conf)) // создание gRPC сервера
	//if err != nil {
	//	log.Fatalf("Failed to listen to port: %s", err)
	//}
	//srv := grpc_server.New(app.Postgres)
	//if err := srv.Grpc.Serve(lis); err != nil {
	//	log.Fatalf("Could not start grpc server: %s", err)
	//}

	//for i := 0; i < amount; i++ {
	//	go generator.Generate(i, telematicCh)
	//}

	//key := 0
	//go func() {
	//	for {
	//		val, ok := <-telematicCh
	//		if ok == false {
	//			close(telematicCh)
	//			break // exit break loop
	//		} else {
	//			val.ID = key
	//			err = newRedis.AddToRedis(val, key) // добавление в редис
	//			if err != nil {
	//				log.Fatalf("Could not add element to redis: %s", err)
	//			}
	//			err = db.AddData(val, ctx) // добавиление в бд
	//			if err != nil && !errors.Is(pgx.ErrNoRows, err) {
	//				log.Fatalf("Could not add element to db: %s", err)
	//			}
	//			err = myKafka.Produce(newKafkaProducer, val) // отправка сообщения в кафку
	//			if err != nil {
	//				log.Fatalf("Could not add element to kafka: %s", err)
	//			}
	//			key++
	//		}
	//	}
	//}()

	// создание grpc сервера
	//var gRPCAddr string
	//flag.StringVar(&gRPCAddr, "grpc-addr", "localhost:8000", "Set the grpc address")
	//
	//lis, err := net.Listen("tcp", ":8000")
	//if err != nil {
	//	log.Fatalf("Failed to listen to port: %s", err)
	//}
	//srv := grpc_server.New(app.Postgres)
	//if err := srv.Grpc.Serve(lis); err != nil {
	//	log.Fatalf("Could not start grpc server: %s", err)
	//}
}
