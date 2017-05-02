package app

import (
	"fmt"
	"gopkg.in/redis.v5"
)

var client *redis.Client

func ConnectToRedis(config *Config) *redis.Client {
	if client != nil {
		return client
	}

	redisConfig := config.RedisConf
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
