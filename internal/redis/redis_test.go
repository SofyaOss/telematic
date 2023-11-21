package redis

import (
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"net"
	"practice/internal/generator"
	"practice/storage"
	"testing"
)

func TestAddToRedis(t *testing.T) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %v", err)
	}

	redisRes, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:         "redis-test",
		Repository:   "redis",
		Tag:          "5-alpine",
		Hostname:     "redis",
		ExposedPorts: []string{"6378"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6378": {
				{HostIP: "0.0.0.0", HostPort: "6378"},
			},
		},
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	defer pool.Purge(redisRes)

	// determine the port the container is listening on
	addr := net.JoinHostPort("localhost", redisRes.GetPort("6379/tcp"))

	err = pool.Retry(func() error {
		client := New(addr)
		for i := 0; i < 7; i++ {
			c := &storage.Car{
				ID:          i,
				Number:      i,
				Speed:       100,
				Coordinates: storage.Coordinates{80, 80},
				Date:        generator.RandomTimestamp(),
			}
			err = client.AddToRedis(c, i)
			if err != nil {
				log.Fatalf("Could not add to redis: %s", err)
			}
		}
		return nil
	})

	if err = pool.Purge(redisRes); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
