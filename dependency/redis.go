package dependency

import (
	"github.com/go-redis/redis/v8"
)

func NewRedis() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0,
	})

	return redis
}
