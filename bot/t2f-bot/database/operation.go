package database

import (
    "banson.moe/t2f-bot/config"
    "context"
    "encoding/json"
    "github.com/go-redis/redis/v8"
    "log"
    "time"
)

var rdb *redis.Client
var ctx = context.Background()
var DefaultTTL = time.Hour * 2

// InitDatabase should be added into main.go
func InitDatabase() *RedisClient {
    rdb = redis.NewClient(&redis.Options{
        Addr:     config.RedisServerAddr,
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

type PictureInfo struct {
    S3Url  string `json:"s3url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

func (client *RedisClient) Get(key string) (*PictureInfo, error) {
    val, err := client.Client.Get(ctx, key).Result()
    if err != nil {
        log.Println("during database Get: key not found")
        return nil, err
    }
    picture := &PictureInfo{}
    err = json.Unmarshal([]byte(val), picture)
    if err != nil {
        log.Printf("during database Get: Unmarshal: %v", err)
        return nil, err
    }
    return picture, nil
}

func (client *RedisClient) Put(key string, pictureInfo *PictureInfo) error {
    pictureInfoBytes, err := json.Marshal(pictureInfo)
    if err != nil {
        log.Printf("during database Put: Marshal: %v", err)
        return err
    }
    err = client.Client.Set(ctx, key, string(pictureInfoBytes), client.TTL).Err()
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
