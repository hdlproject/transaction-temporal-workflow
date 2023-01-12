package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Idempotency struct {
	redis *redis.Client
}

func NewIdempotency(redis *redis.Client) Idempotency {
	return Idempotency{
		redis: redis,
	}
}

func (i Idempotency) IsAllowed(idempotencyKey string) (bool, error) {
	res := i.redis.SetNX(context.Background(), idempotencyKey, "", time.Hour)
	if res.Err() != nil {
		return false, res.Err()
	}

	return res.Val(), nil
}
