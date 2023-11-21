package generator

import (
	"fmt"
	"log"
	"practice/storage"

	//"github.com/golang/protobuf/protoc-gen-go/generator"
	"math/rand"
	"testing"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestGenerator(t *testing.T) {
	kafkaCh := make(chan *storage.Car)
	log.Println("Starting generator...")
	go Generate(rand.Intn(3), kafkaCh)

	for {
		val, ok := <-kafkaCh
		if ok == false {
			break // exit break loop
		} else {
			fmt.Println(val, ok)
		}
	}
}
