package app

import (
	"fmt"
	"gopkg.in/redis.v5"
	"testing"
)

func init() {
	ConnectToRedis(LoadConfig())
}

func TestConnectToRedis(t *testing.T) {
	client := GetRedisClient()

	defer func(client *redis.Client) {
		fmt.Println("Closing Redis client...")
		err := client.Close()
		if err != nil {
			fmt.Printf("%v\n", "Closing Redis client failed!")
			panic(err)
		}
	}(client)

	fmt.Println(client.Info())
}
