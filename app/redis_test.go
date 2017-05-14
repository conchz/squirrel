package app_test

import (
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
			fmt.Printf("%v\n", "Closing Redis client failed!")
			panic(err)
		}
	}(client)

	fmt.Println(client.Info())
}
