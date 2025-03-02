package redis

import (
	"context"
	"fmt"
	"os"

	redis "github.com/go-redis/redis/v8"
)

var Client *redis.Client

func ConnectToRedis() error {
	ctx := context.Background()
	// Retrieve Redis address from environment variables
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		return fmt.Errorf("REDIS_URL environment variable is not set")
	}
	// Create a new Redis client
	Client = redis.NewClient(&redis.Options{
		Addr:     redisAddr, // e.g., "localhost:6379"
		Password: "",        // No password
		DB:       0,         // Default DB
	})
	// Test the connection
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}
