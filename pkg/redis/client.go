package client

import (
	"bank_micro/pkg/config"
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.CoreCfg.RedisAddr,
		DB:       0,  // Default DB
		PoolSize: 10, // Connection pool size
	})

	// Test connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}
}

func SetCache(key string, value interface{}, ttl time.Duration) error {
	return RedisClient.Set(context.Background(), key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return RedisClient.Get(context.Background(), key).Result()
}
