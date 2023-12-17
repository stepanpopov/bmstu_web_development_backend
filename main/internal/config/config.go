package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

const secret = "rust manipulates memory much more efficiently"

type Config struct {
	ServiceHost       string
	ServicePort       int
	CalculateSecret   string
	CalculateCallback string
	JWT               JWTConfig
	Redis             RedisConfig
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

func NewConfig(ctx context.Context) (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	// viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.JWT.Secret = secret
	cfg.JWT.ExpiresIn = time.Duration(time.Hour * 24 * 365)

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))

	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}

	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)

	cfg.CalculateSecret = "rust"
	cfg.CalculateCallback = "http://localhost:8765/calculate/"

	return cfg, nil
}
