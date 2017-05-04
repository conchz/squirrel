package app

import "gopkg.in/redis.v5"

var client *redis.Client

func init() {
	connectToRedis(LoadConfig())
}

func connectToRedis(config *Config) {
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
