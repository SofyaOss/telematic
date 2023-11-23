package kafka

import (
	"encoding/json"
	"log"

	"practice/storage"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	prod  *kafka.Producer
	topic string
}

type Consumer struct {
	cons *kafka.Consumer
}

func NewProducer(addr string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": addr, // "kafka:9092"
	}
	topic := "telematicTopic"
	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("could not connect to kafka: %s", err)
		return nil, err
	}
	p := Producer{prod: producer, topic: topic}
	return &p, nil
}

func Produce(producer *Producer, item *storage.Car) error {
	mes, err := json.Marshal(item)
	if err != nil {
		log.Fatalf("Could not convert data to json: %s", err)
	}
	err = producer.prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &producer.topic, Partition: kafka.PartitionAny},
		Value:          mes,
	}, nil)
	if err != nil {
		log.Println("кафка блять", err)
		return err
	}
	return nil
	//else {
	//	log.Println("победа")
	//}
	//log.Println(mes, ok)
}
