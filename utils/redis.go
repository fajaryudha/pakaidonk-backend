package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // default Redis
})

var ctx = context.Background()

func SaveTokenToRedis(jti string, token string, ttl time.Duration) error {
	return redisClient.Set(ctx, jti, token, ttl).Err()
}
