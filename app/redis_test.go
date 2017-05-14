package app

import (
	"gopkg.in/redis.v5"
	"testing"
)

func init() {
	ConnectToRedis(LoadConfig())
}

func TestConnectToRedis(t *testing.T) {
	client := GetRedisClient()

	defer func(client *redis.Client) {
		t.Log("Closing Redis client...")
		err := client.Close()
		if err != nil {
			t.Errorf("%v\n", "Closing Redis client failed!")
			panic(err)
		}
	}(client)

	t.Log(client.Info())
}
