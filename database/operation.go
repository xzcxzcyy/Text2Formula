package database

import (
	"context"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()
var DefaultTTL = time.Hour * 2

// initDatabase should be added into main.go
func initDatabase() *RedisClient {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{
		Client: rdb,
		TTL:    DefaultTTL,
	}
}

type RedisClient struct {
	Client *redis.Client
	TTL    time.Duration
}

type Operation interface {
	Get(key string) (string, error)

	Put(key, url string) error

	Update(key, url string) error

	SetTTL(hour uint64) error

	Delete(key string) error

}


func (client *RedisClient) Get(key string) (string, error) {
	val, err := client.Client.Get(ctx, key).Result()
	if err != nil {
		log.Println("key not found")
		return "", err
	}
	return val, nil
}

func (client *RedisClient) Put(key string, url string) error {
	err := client.Client.Set(ctx, key, url, client.TTL).Err()
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
	err := client.Client.Del(ctx, key).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (client *RedisClient) Update(key, url string) error {
	client.Client.Set(ctx, key, url, client.TTL)
	return nil
}
