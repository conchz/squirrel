package app

import (
	"fmt"
	"gopkg.in/redis.v6"
)

var client *redis.Client

func ConnectToRedis() *redis.Client {
	if client != nil {
		return client
	}

	redisConfig := GetConfig().RedisConf
	client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       0,
	})

	result, err := client.Ping().Result()
	if err != nil {
		fmt.Println(result, err)
		panic(err)
	}

	return client
}
