package utils

import (
	"github.com/go-redis/redis"
)

// CreateRedisClient returns a new redis client.
func CreateRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return client
}
