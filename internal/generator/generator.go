package generator

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"math/rand"
	"practice/storage"
	"time"
)

type telematic struct {
	timestamp time.Time
	speed     int
	coords    *storage.Coordinates
}

func RandomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
	randomNow := time.Unix(randomTime, 0)
	return randomNow
}

func Generate(num int, c chan *storage.Car) {
	//var tList []*telematic
	prevTimestamp := RandomTimestamp()
	log.Println("time is", prevTimestamp)
	prevCoords := storage.Coordinates{-90 + rand.Float64()*180, -180 + rand.Float64()*360} //res[i] = min + rand.Float64() * (max - min)
	for i := 0; i < 5; i++ {
		newTimestamp := prevTimestamp.Add(time.Duration(rand.Intn(60)) * time.Second)
		newSpeed := rand.Intn(120)
		//fmt.Println(newTimestamp.Unix(), newTimestamp.Unix()/3600)
		//r := (newTimestamp.Unix()/3600 - prevTimestamp.Unix()/3600) * int64(newSpeed)
		//fmt.Println(prevTimestamp, newTimestamp, r)
		var newCoords storage.Coordinates
		//x := -float64(r) + rand.Float64()*(float64(r)*2)
		//y := math.Pow(float64(r), 2) - math.Pow(x, 2)
		//fmt.Println(r, x, y)
		newCoords = storage.Coordinates{prevCoords.Latitude, prevCoords.Longitude} // это заглушка, пересчитать все
		//if r < 1 {
		//	newCoords = &coordinates{prevCoords.latitude + (-90 + rand.Float64()*180), -180 + rand.Float64()*360}
		//} else {
		//	newCoords = &coordinates{-90 + rand.Float64()*180, prevCoords.longitude + (-180 + rand.Float64()*360)}
		//}
		//fmt.Printf("prevCoords: %v, prevTime: %v, newCoords: %v, newTime: %v\n", prevCoords, prevTimestamp, newCoords, newTimestamp)
		//t := &telematic{
		//	timestamp: newTimestamp,
		//	speed:     newSpeed,
		//	coords:    newCoords,
		//}
		car := &storage.Car{0, num, newSpeed, newCoords, newTimestamp}
		c <- car
		//tList = append(tList, t)
		prevTimestamp = newTimestamp
		prevCoords = newCoords
		time.Sleep(1 * time.Second)
	}
	return
	//return tList
}

func Generator() {
	kafkaCh := make(chan *storage.Car)
	fmt.Println("created")
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}
	topic := "telematicTopic"
	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(err)
	}
	go Generate(1, kafkaCh)
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
	//for i := range kafkaCh {
	//	fmt.Println(i)
	//}
	//return nil
}
