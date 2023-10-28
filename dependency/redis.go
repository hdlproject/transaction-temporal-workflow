package dependency

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddress"),
		Password: os.Getenv("RedisPassword"),
		DB:       0,
	})

	return redisClient
}
