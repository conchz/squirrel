package app

import (
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
	"testing"
)

func TestConnectToRedis(t *testing.T) {
	config := LoadConfig()
	client := ConnectToRedis(config)

	defer closeRedisClient(client)

	fmt.Println(client.Info())
}

func closeRedisClient(client *redis.Client) {
	if client != nil {
		fmt.Println("Closing Redis client...")
		err := client.Close()
		if err != nil {
			panic(errors.New("Closing Redis client failed!"))
		}
	} else {
		panic(errors.New("There is no available Redis client!"))
	}
}
