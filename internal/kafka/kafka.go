package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"practice/storage"
)

type Producer struct {
	prod  *kafka.Producer
	topic string
}

type Consumer struct {
	cons *kafka.Consumer
}

func NewProducer() (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
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
		log.Println("AAAAAAAAAAAAA kafkaaaaa", err)
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
