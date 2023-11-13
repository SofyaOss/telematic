package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"math/rand"
	"practice/internal/generator"
	"practice/storage"
	"testing"
)

func TestProduce(t *testing.T) {
	newKafkaProd := NewProducer()
	for i := 1; i < 9; i++ {
		newCar := storage.Car{
			i,
			rand.Intn(3),
			100,
			storage.Coordinates{80, 80},
			generator.RandomTimestamp(),
		}
		mes, err := json.Marshal(newCar)
		if err != nil {
			log.Println("AAAAAAAAAAAAA kafkaaaaa", err)
		}
		err = newKafkaProd.prod.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &newKafkaProd.topic, Partition: kafka.PartitionAny},
			Value:          mes,
		}, nil)
		if err != nil {
			log.Println("кафка блять", err)
		}
	}
}
