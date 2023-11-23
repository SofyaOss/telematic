package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"practice/internal/grpc_server"
	"strconv"

	"practice/internal/config"
	"practice/internal/generator"
	myKafka "practice/internal/kafka"

	"github.com/jackc/pgx/v4"
)

func (a *Application) refreshTelematicTable(ctx context.Context) {
	err := a.Postgres.DropTable(ctx) // удаление старой таблицы
	if err != nil {
		log.Println("Could not drop the table:", err)
	}

	err = a.Postgres.CreateTable(ctx) // создание новой таблицы
	if err != nil {
		log.Println("Could not create the table", err)
	}
}

func (a *Application) startApp(ctx context.Context) {
	amount, err := strconv.Atoi(os.Getenv("TRANSPORT_AMOUNT")) // количество машин
	if err != nil {
		log.Fatal("Transport amount must be an integer")
	}

	for i := 0; i < amount; i++ {
		go generator.GenerateTelematic(i, a.Channel)
	}

	key := 0
	go func() {
		for {
			val, ok := <-a.Channel
			if ok == false {
				close(a.Channel)
				break // exit break loop
			} else {
				val.ID = key
				err := a.Redis.AddToRedis(val, key) // добавление в редис
				if err != nil {
					log.Fatalf("Could not add element to redis: %s", err)
				}
				err = a.Postgres.AddData(val, ctx) // добавиление в бд
				if err != nil && !errors.Is(pgx.ErrNoRows, err) {
					log.Fatalf("Could not add element to db: %s", err)
				}
				err = myKafka.Produce(a.Kafka, val) // отправка сообщения в кафку
				if err != nil {
					log.Fatalf("Could not add element to kafka: %s", err)
				}
				key++
			}
		}
	}()
}

func (a *Application) startServer(addr string) {
	lis, err := net.Listen("tcp", addr) // создание gRPC сервера
	if err != nil {
		log.Fatalf("Failed to listen to port: %s", err)
	}
	srv := grpc_server.New(a.Postgres)
	if err := srv.Grpc.Serve(lis); err != nil {
		log.Fatalf("Could not start grpc server: %s", err)
	}
}

func getConfig() *config.AppConf {
	var conf *config.AppConf
	conf = config.ReadConf(conf)
	return conf
}
