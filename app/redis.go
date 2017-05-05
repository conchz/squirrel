package app

import "gopkg.in/redis.v5"

var client *redis.Client

func ConnectToRedis(config *Config) {
	redisConfig := config.RedisConf
	client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return client
}
