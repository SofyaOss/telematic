package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"practice/internal/generator"
	"practice/storage"
	"practice/storage/postgres"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %v", err)
	}

	opts := dockertest.RunOptions{
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

	log.Println("1")
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not start resourse: %s", err)
	}

	log.Println("2")
	var db *postgres.TelematicDB
	if err = pool.Retry(func() error {
		dbConnStr := fmt.Sprintf("host=localhost port=5432 user=postgres dbname=testdb password=postgres sslmode=disable")
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		db, err = postgres.New(ctx, dbConnStr) // подключение к базе данных
		if err != nil {
			log.Println("db not ready yet")
			return err
		}
		return db.DropTable()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	//err = pool.Client.Ping()
	//if err != nil {
	//	log.Fatalf("could not connect to docker: %v", err)
	//}

	//res, err := db.GetAllData()
	//if err != nil {
	//	log.Println("failed to get all data with error", err)
	//} else {
	//	log.Println("GetAllData() result:", res)
	//}
	err = db.DropTable()
	if err != nil {
		log.Println("failed to drop table", err)
	}

	err = db.CreateTable()
	if err != nil {
		log.Println("failed to create table", err)
	}

	log.Println("trying to add data to test db...")
	kafkaCh := make(chan *storage.Car) // создание канала для передачи телематики
	go generator.Generate(1, kafkaCh)
	for i := 0; i < 8; i++ {
		val, ok := <-kafkaCh
		if ok == false {
			log.Println(val, ok, "<-- loop broke!")
			close(kafkaCh)
			break // exit break loop
		} else {
			err = db.AddData(val)
			if err != nil && !errors.Is(pgx.ErrNoRows, err) {
				log.Fatalf("could not add element to db: %s", err)
			}
		}
	}
	code := m.Run()
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(code)
}
