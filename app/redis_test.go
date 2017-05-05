package app_test

import (
	"errors"
	"fmt"
	"github.com/lavenderx/squirrel/app"
	"gopkg.in/redis.v5"
	"testing"
)

func TestConnectToRedis(t *testing.T) {
	app.ConnectToRedis(app.LoadConfig())

	client := app.GetRedisClient()

	defer func(client *redis.Client) {
		fmt.Println("Closing Redis client...")
		err := client.Close()
		if err != nil {
			panic(errors.New("Closing Redis client failed!"))
		}
	}(client)

	fmt.Println(client.Info())
}
