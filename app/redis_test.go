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

	defer func(client *redis.Client) {
		fmt.Println("Closing Redis client...")
		err := client.Close()
		if err != nil {
			panic(errors.New("Closing Redis client failed!"))
		}
	}(client)

	fmt.Println(client.Info())
}
