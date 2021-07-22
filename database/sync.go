package database

import "github.com/pkg/errors"

var (
	ClientNotSetUp error = errors.New("client has not set up")
)

// TODO: check the format of tex, like: if the text should have $...$ or not if the tex has redundent space... need a filter function
func SetTTL(client *RedisClient, TTL uint64) error {
	if client == nil {
		return ClientNotSetUp
	}
	return client.SetTTL(TTL)
}

func UploadPicture(client *RedisClient, tex string, CloudUrl string) error {
	if client == nil {
		return ClientNotSetUp
	}
	return client.Put(tex, CloudUrl)
}

func isOnCloud(client *RedisClient, tex string) (CloudUrl string, err error) {
	if client == nil {
		return "", ClientNotSetUp
	}
	val, err := client.Get(tex)
	if err != nil {
		return "", err
	}
	return val, err
}
