package main

import (
	"context"
	"flag"
	"practice/internal/grpc_server"
	"strconv"

	//"encoding/json"
	//"fmt"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"net/http"
	"os"
	"practice/internal/generator"
	myKafka "practice/internal/kafka"
	"practice/internal/redis"
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

	//res, err := db.GetAllData()
	//if err != nil {
	//	log.Println("да блять")
	//}
	//log.Println("GetAllData() result:", res)
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

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "redis:6379",
	//	Password: "",
	//	DB:       0,
	//})

	newRedis := redis.New()
	newKafkaProducer, err := myKafka.NewProducer()
	if err != nil {
		log.Fatalf("could not connect to kafka: %s", err)
	}

	kafkaCh := make(chan *storage.Car)
	amount, err := strconv.Atoi(os.Getenv("TRANSPORT_AMOUNT"))
	if err != nil {
		log.Fatal("Transport amount must be an integer")
	}
	for i := 0; i < amount; i++ {
		go generator.Generate(i, kafkaCh)
	}

	key := 0
	for {
		val, ok := <-kafkaCh
		if ok == false {
			log.Println(val, ok, "<-- loop broke!")
			close(kafkaCh)
			break // exit break loop
		} else {
			err = redis.AddToRedis(newRedis, val, key)
			if err != nil {
				log.Fatalf("could not add element to redis: %s", err)
			}
			err = myKafka.Produce(newKafkaProducer, val)
			if err != nil {
				log.Fatalf("could not add element to kafka: %s", err)
			}
			key++
		}
		if key >= 24 {
			break
		}
	}

	//key := 0
	//for {
	//	val, ok := <-kafkaCh
	//	if ok == false {
	//		log.Println(val, ok, "<-- loop broke!")
	//		close(kafkaCh)
	//		break // exit break loop
	//	} else {
	//		mes, err := json.Marshal(val)
	//		if err != nil {
	//			log.Println("AAAAAAAAAAAAA", err)
	//		}
	//
	//		if key < 1000 {
	//			err = client.Set(strconv.Itoa(key), mes, 0).Err()
	//			if err != nil {
	//				log.Println("aaaaaaaaa")
	//			} else {
	//				key++
	//			}
	//		} else {
	//			client.Del(strconv.Itoa(key - 1000))
	//			err = client.Set(strconv.Itoa(key), mes, 0).Err()
	//			if err != nil {
	//				log.Println("aaaaaaaaa2")
	//			} else {
	//				key++
	//			}
	//		}
	//
	//	}
	//}

	//redisVar, err := client.Get("2").Result()
	//if err != nil {
	//	log.Println("проверка", err)
	//} else {
	//	log.Println("победа", redisVar)
	//}
	//log.Println("закончили упражнение")
	//client.Del("2")
	//redisVar, err = client.Get("2").Result()
	//if err != nil {
	//	log.Println("допустим", err)
	//}
	//log.Println(redisVar)

	//fmt.Println("created")
	//config := &kafka.ConfigMap{
	//	"bootstrap.servers": "kafka:9092",
	//}
	//topic := "telematicTopic"
	//producer, err := kafka.NewProducer(config)
	//if err != nil {
	//	panic(err)
	//}

	//key := 0
	//for {
	//	val, ok := <-kafkaCh
	//	if ok == false {
	//		log.Println(val, ok, "<-- loop broke!")
	//		break // exit break loop
	//	} else {
	//		mes, err := json.Marshal(val)
	//		if err != nil {
	//			log.Println("AAAAAAAAAAAAA", err)
	//		}
	//		err = producer.Produce(&kafka.Message{
	//			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//			Value:          mes,
	//		}, nil)
	//		if err != nil {
	//			log.Println("кафка блять", err)
	//		}
	//else {
	//	log.Println("победа")
	//}
	//log.Println(mes, ok)
	//	}
	//	key++
	//	if key >= 24 {
	//		break
	//	}
	//
	//}
	//log.Println("were here")
	//producer.Flush(15 * 1000)
	//producer.Close()

	//log.Println("start consumer")
	//time.Sleep(5 * time.Second)
	//log.Println("start consumer x2")
	//
	//config2 := &kafka.ConfigMap{
	//	"bootstrap.servers": "kafka:9092",
	//	"group.id":          "myGroup",
	//	"auto.offset.reset": "earliest",
	//}
	//consumer, err := kafka.NewConsumer(config2)
	//if err != nil {
	//	panic(err)
	//}
	//consumer.SubscribeTopics([]string{"telematicTopic"}, nil)
	//for {
	//	msg, err := consumer.ReadMessage(-1)
	//	if err == nil {
	//		log.Printf("recieved message: %s\n", string(msg.Value))
	//	} else {
	//		log.Printf("Error while consuming message %v(%v)\n", err, msg)
	//	}
	//}
	//consumer.Close()

	//router := mux.NewRouter()
	//router.HandleFunc("/data", db.GetAllData()).Methods("GET")
	//router.HandleFunc("/add", postgres.AddData(db)).Methods("POST")

	//log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))

	//lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	//if err != nil {
	//	log.Fatalf("failed to listen: %s", err)
	//}
	//serv := grpc.NewServer()
	//pb.RegisterGRPCServiceServer(serv, &grpc_server.Server{})
	//if err = serv.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %s", err)
	//}
	var gRPCAddr string

	flag.StringVar(&gRPCAddr, "grpc-addr", "localhost:11000", "Set the grpc address")
	srv := grpc_server.New(gRPCAddr, db)
	if err := srv.Open(); err != nil {
		log.Fatalf("failed to start server: %s", err)
	} else {
		log.Println("урааа")
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
