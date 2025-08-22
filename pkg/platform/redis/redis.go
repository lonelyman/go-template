package platform

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// InitRedis initializes default redis connection
func InitRedis() (*redis.Client, error) {
	return InitRedisWithName("redis")
}

// InitRedisWithName initializes redis connection with custom configuration name
func InitRedisWithName(connectionName string) (*redis.Client, error) {
	envPrefix := strings.ToUpper(connectionName)

	// Get configuration values
	host := viper.GetString(fmt.Sprintf("%s.host", connectionName))
	port := viper.GetInt(fmt.Sprintf("%s.port", connectionName))
	password := viper.GetString(fmt.Sprintf("%s.password", connectionName))
	db := viper.GetInt(fmt.Sprintf("%s.db", connectionName))
	maxRetries := viper.GetInt(fmt.Sprintf("%s.max_retries", connectionName))
	poolSize := viper.GetInt(fmt.Sprintf("%s.pool_size", connectionName))
	minIdleConns := viper.GetInt(fmt.Sprintf("%s.min_idle_conns", connectionName))

	// Set default values if not provided
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 6379
	}
	if maxRetries == 0 {
		maxRetries = 3
	}
	if poolSize == 0 {
		poolSize = 10
	}
	if minIdleConns == 0 {
		minIdleConns = 2
	}

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		MaxRetries:   maxRetries,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolTimeout:  10 * time.Second,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis (%s): %w", envPrefix, err)
	}

	log.Printf("Successfully connected to Redis (%s) at %s:%d", envPrefix, host, port)

	return rdb, nil
}
