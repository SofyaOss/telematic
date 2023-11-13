package generator

import (
	"fmt"
	"practice/storage"

	//"github.com/golang/protobuf/protoc-gen-go/generator"
	"math/rand"
	"testing"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestGenerator(t *testing.T) {
	kafkaCh := make(chan *storage.Car)
	fmt.Println("created")
	go Generate(rand.Intn(3), kafkaCh)

	//config := &kafka.ConfigMap{
	//	"bootstrap.servers": "localhost:9092",
	//}
	//topic := "testTopic"
	//producer, err := kafka.NewProducer(config)
	//if err != nil {
	//	panic(err)
	//}

	for {
		val, ok := <-kafkaCh
		if ok == false {
			fmt.Println(val, ok, "<-- loop broke!")
			break // exit break loop
		} else {
			fmt.Println(val, ok)
		}
	}

}
