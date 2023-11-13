package redis

import (
	"encoding/json"
	"log"
	"math/rand"
	"practice/internal/generator"
	"practice/storage"
	"strconv"
	"testing"
)

func TestAddToRedis(t *testing.T) {
	newClient := New()
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
			log.Println("AAAAAAAAAAAAA", err)
		}
		err = newClient.client.Set(strconv.Itoa(i), mes, 0).Err()
		if err != nil {
			log.Println("aaaaaaaaa")
		} else {
			log.Println("telematic added")
		}
	}
}
