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

func New() *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	c := Client{client: client}
	return &c
}

func AddToRedis(client *Client, c *storage.Car, key int) error {
	mes, err := json.Marshal(c)
	if err != nil {
		log.Println("AAAAAAAAAAAAA", err)
		return err
	}

	if key < 20 {
		err = client.client.Set(strconv.Itoa(key), mes, 0).Err()
		if err != nil {
			log.Println("aaaaaaaaa")
			return err
		} else {
			key++
			log.Println("< 1000")
		}
	} else {
		client.client.Del(strconv.Itoa(key - 20))
		err = client.client.Set(strconv.Itoa(key), mes, 0).Err()
		if err != nil {
			log.Println("aaaaaaaaa2")
			return err
		} else {
			key++
			log.Println("> 1000")
		}
	}
	return nil
}
