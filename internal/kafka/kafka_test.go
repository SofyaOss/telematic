package kafka

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"practice/internal/generator"
	"practice/storage"
	"testing"
)

func TestProduce(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	zookeeperResource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "zookeeper-test",
		Repository:   "confluentinc/cp-zookeeper",
		Tag:          "latest",
		Hostname:     "zookeeper",
		ExposedPorts: []string{"2182"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"2182": {
				{HostIP: "0.0.0.0", HostPort: "2182"},
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not start zookeeper: %s", err)
	}
	defer pool.Purge(zookeeperResource)

	kafkaRes, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "kafka-test",
		Repository: "confluentinc/cp-kafka",
		Tag:        "latest",
		Hostname:   "kafka",
		Env: []string{
			"KAFKA_BROKER_ID: 1",
			"KAFKA_ZOOKEEPER_CONNECT: zookeeper:2182",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1",
			"KAFKA_ADVERTISED_LISTENERS: INSIDE://:9092,OUTSIDE://:9093",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INSIDE:PLAINTEXT,OUTSIDE:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME: INSIDE",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"9093/tcp": {{HostIP: "localhost", HostPort: "9093/tcp"}},
		},
		ExposedPorts: []string{"9093/tcp"},
	})
	if err != nil {
		log.Fatalf("Could not start kafka: %s", err)
	}
	defer pool.Purge(kafkaRes)

	err = pool.Retry(func() error {
		newKafkaProducer, err := NewProducer() // создание продюсера кафки
		if err != nil {
			log.Fatalf("Could not connect to kafka: %s", err)
		}
		for i := 0; i < 7; i++ {
			c := &storage.Car{
				ID:          i,
				Number:      i,
				Speed:       100,
				Coordinates: storage.Coordinates{80, 80},
				Date:        generator.RandomTimestamp(),
			}
			err = Produce(newKafkaProducer, c)
			if err != nil {
				log.Fatalf("Could not send message to kafka: %s", err)
			}
		}
		return nil
	})

	if err = pool.Purge(kafkaRes); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err = pool.Purge(zookeeperResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
