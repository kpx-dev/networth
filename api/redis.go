package main

import (
	"github.com/go-redis/redis"
)

// RedisClient redis client struct
type RedisClient struct {
	*redis.Client
}

// NewRedisClient new redis client
func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_HOST", "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{client}
}
