package redisconn

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	if os.Getenv("ENV") == "production" {
		redisURL := os.Getenv("REDIS_URL")
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("failed to parse redis url: %v", err)
		}
		return redis.NewClient(opt)
	} else {
		// Local development
		return redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		})
	}
}
