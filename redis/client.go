package redisconn

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	var redisURL string
	if os.Getenv("ENV") == "production" {
		redisURL = os.Getenv("REDIS_URL")
	} else {
		redisURL = "localhost:6379"
	}
	return redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})
}
