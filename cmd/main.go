package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"practice/internal/generator"
	"practice/storage"
	"practice/storage/postgres"
	"time"
)

func main() {
	//connect to database
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()

	//create the table if it doesn't exist
	//log.Println("im working")
	//_, err = db.Exec("DROP TABLE IF EXISTS telematic")

	//_, err = db.Exec("CREATE TABLE IF NOT EXISTS telematic (id SERIAL PRIMARY KEY, car_number INT, speed INT, latitude TEXT, longitude TEXT, date TEXT)")
	// , pub_date BIGINT DEFAULT 0)
	//if err != nil {
	//	log.Fatal(err)
	//}
	log.Println("ну с богом...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := postgres.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("бля...", err)
	}

	err = db.DropTable()
	if err != nil {
		log.Println("да бля", err)
	}

	err = db.CreateTable()
	if err != nil {
		log.Println("бляяяяя", err)
	}

	res, err := db.GetAllData()
	if err != nil {
		log.Println("да блять")
	}
	log.Println("GetAllData() result:", res)
	/*
		for i := 1; i < 9; i++ {
			newCar := storage.Car{
				i,
				rand.Intn(3),
				100,
				generator.Coordinates{80, 80},
				generator.RandomTimestamp(),
			}
			err = db.AddData(newCar)
			if err != nil {
				log.Println("так блять", err)
			} else {
				log.Println("победа")
			}
		}

		res1, err := db.GetAllData()
		if err != nil {
			log.Println("ладно блять", err)
		} else {
			log.Println("GetAllData() result1:", res1)
		}

		res2, err := db.GetByDate("1995-01-01", "2020-01-01")
		if err != nil {
			log.Println("да еб", err)
		} else {
			log.Println("победа х2", res2)
		}

		var nums []int
		nums = append(nums, 2)
		res3, err := db.GetByCarNumber(nums)
		if err != nil {
			log.Println("да ебать тебя", err)
		} else {
			log.Println("победа х3", res3)
		}
	*/

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	kafkaCh := make(chan *storage.Car)
	for i := 0; i < 3; i++ {
		go generator.Generate(i, kafkaCh)
	}

	for {
		val, ok := <-kafkaCh
		if ok == false {
			log.Println(val, ok, "<-- loop broke!")
			close(kafkaCh)
			break // exit break loop
		} else {
			mes, err := json.Marshal(val)
			if err != nil {
				log.Println("AAAAAAAAAAAAA", err)
			}

			err = client.Set("1", mes, 0).Err()
			if err != nil {
				log.Println("aaaaaaaaa")
			}

			redisVar, err := client.Get("1").Result()
			if err != nil {
				log.Println("блять редис ты то че начинаешь", err)
			}
			log.Println("победа", redisVar)
		}
	}
	/*
		fmt.Println("created")
		config := &kafka.ConfigMap{
			"bootstrap.servers": "localhost:9092",
		}
		topic := "telematicTopic"
		producer, err := kafka.NewProducer(config)
		if err != nil {
			panic(err)
		}

		for {
			val, ok := <-kafkaCh
			if ok == false {
				log.Println(val, ok, "<-- loop broke!")
				break // exit break loop
			} else {
				mes, err := json.Marshal(val)
				if err != nil {
					log.Println("AAAAAAAAAAAAA", err)
				}
				err = producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          mes,
				}, nil)
				if err != nil {
					log.Println("кафка блять", err)
				} else {
					log.Println("победа")
				}
				//log.Println(mes, ok)
			}
		}
		log.Println("were here")
		producer.Flush(15 * 1000)
		producer.Close()
	*/

	router := mux.NewRouter()
	//router.HandleFunc("/data", db.GetAllData()).Methods("GET")
	//router.HandleFunc("/add", postgres.AddData(db)).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
