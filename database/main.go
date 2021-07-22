package main

import (
	"context"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()
var DefaultTTL = time.Hour * 2

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

func main() {
	client := initDatabase()
	client.Put("hello", "world")
	log.Println(client.Get("hello"))
	client.Delete("hello")
	log.Println(client.Get("hello"))
}