package main

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type RedisClient struct {
	Client *redis.Client
	TTL    time.Duration
}

type Operation interface {
	Get(key string) (string, error)

	Put(key, url string) error

	SetTTL(hour uint64) error

	Delete(key string) error

}


func (client *RedisClient) Get(key string) (string, error) {
	val, err := client.Client.Get(key).Result()
	if err != nil {
		log.Println("key not found")
		return "", err
	}
	return val, nil
}

func (client *RedisClient) Put(key string, url string) error {
	err := client.Client.Set(key, url, client.TTL).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (client *RedisClient) SetTTL(hour uint64) error {
	client.TTL = time.Duration(hour) * time.Hour
	return nil
}

func (client *RedisClient) Delete(key string) error {
	err := client.Client.Del(key).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
