package redis

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"mygo/template/pkg/config"
)

var (
	defaultRedisClient *redis.Client

	defaultRedisClientOnce sync.Once
)

func InitRedisClient(redisConfig map[string]*config.Redis) {
	defaultRedisConfig, ok := redisConfig["default"]
	if !ok {
		panic("redis default should be configured")
	}
	initDefaultRedisClient(defaultRedisConfig)
}

func initDefaultRedisClient(redisConfig *config.Redis) {
	if defaultRedisClient != nil {
		return
	}

	defaultRedisClientOnce.Do(func() {
		var err error
		defaultRedisClient, err = newRedisClient(redisConfig)
		if err != nil {
			panic(fmt.Errorf("failed to connect to redis, err: %v", err))
		}
	})
}

func newRedisClient(redisConfig *config.Redis) (*redis.Client, error) {
	opt := &redis.Options{
		Addr:         redisConfig.Addr,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,

		PoolSize:     20 * runtime.NumCPU(),
		MinIdleConns: 10 * runtime.NumCPU(),
		IdleTimeout:  180 * time.Second,
	}

	if redisConfig.DialTimeout > 0 {
		opt.DialTimeout = time.Duration(redisConfig.DialTimeout) * time.Second
	}
	if redisConfig.ReadTimeout > 0 {
		opt.ReadTimeout = time.Duration(redisConfig.ReadTimeout) * time.Second
	}
	if redisConfig.WriteTimeout > 0 {
		opt.WriteTimeout = time.Duration(redisConfig.WriteTimeout) * time.Second
	}

	if redisConfig.PoolSize > 0 {
		opt.PoolSize = redisConfig.PoolSize
	}
	if redisConfig.MinIdleConns > 0 {
		opt.MinIdleConns = redisConfig.MinIdleConns
	}
	if redisConfig.IdleTimeout > 0 {
		opt.IdleTimeout = time.Duration(redisConfig.IdleTimeout) * time.Second
	}

	client := redis.NewClient(opt)
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}

func GetDefaultRedisClient() *redis.Client {
	return defaultRedisClient
}
