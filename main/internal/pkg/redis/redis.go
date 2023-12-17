package redis

import (
	"context"
	"fmt"
	"rip/internal/config"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const servicePrefix = "rip_service." // наш префикс сервиса

type RedisClient struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*RedisClient, error) {
	client := &RedisClient{}

	client.cfg = cfg

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Username:    cfg.User,
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	client.client = redisClient

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return client, nil
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
