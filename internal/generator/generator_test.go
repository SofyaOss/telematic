package generator

import (
	"fmt"
	"testing"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
)

func TestGenerator(t *testing.T) {
	kafkaCh := make(chan *telematic)
	fmt.Println("created")
	go generate(kafkaCh)

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
			fmt.Println(val.timestamp, ok)
		}
	}

}
