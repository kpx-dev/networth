package main

import (
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/networth-app/networth/api/lib"
)

// RedisClient redis client struct
type RedisClient struct {
	*redis.Client
}

// NewRedisClient new redis client
func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     nwlib.GetEnv("REDIS_HOST", "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisClient{client}
}

// GetNetworth get current networth
func (c *RedisClient) GetNetworth() float64 {
	val, err := c.Get("networth").Result()

	if err != nil {
		log.Println(err)
	}

	networth, _ := strconv.ParseFloat(val, 64)

	return networth
}
