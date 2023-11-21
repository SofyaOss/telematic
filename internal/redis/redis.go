package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"practice/storage"
	"strconv"
)

type Client struct {
	client *redis.Client
}

func New(addr string) *Client { // создание нового клиента редис
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	c := Client{client: client}
	return &c
}

func (c *Client) AddToRedis(car *storage.Car, key int) error { // отправка сообщения в редис
	mes, err := json.Marshal(car)
	if err != nil {
		return err
	}

	if key < 1000 {
		err = c.client.Set(strconv.Itoa(key), mes, 0).Err()
		if err != nil {
			log.Fatalf("Could not set value in redis: %s", err)
			return err
		} else {
			key++
		}
	} else {
		c.client.Del(strconv.Itoa(key - 1000))
		err = c.client.Set(strconv.Itoa(key), mes, 0).Err()
		if err != nil {
			log.Fatalf("Could not set value in redis: %s", err)
			return err
		} else {
			key++
		}
	}
	return nil
}
