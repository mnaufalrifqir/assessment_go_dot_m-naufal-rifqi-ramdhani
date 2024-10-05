package database

import (
	"api-dot/utils"
	"context"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	Client = redis.NewClient(&redis.Options{
		Addr:     utils.GetConfig("REDIS_HOST") + ":" + utils.GetConfig("REDIS_PORT"),
		Password: utils.GetConfig("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
