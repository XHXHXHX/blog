package redis

import (
	redisConfig "blog/config/redis"
	"blog/library/log"
	"errors"
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

func GetNewClient() (*redis.Client, error) {
	config, err := redisConfig.GetRedisConfig()
	if err != nil {
		log.New().Error("Redis config write failed")
		return nil, errors.New("Redis config write failed")
	}

	conn := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:    string(uc.Con.Redis.WritePassword),
		DB:          0,
		DialTimeout: time.Second * time.Duration(uc.Con.Redis.DialTimeout),
	})

	return conn, nil
}