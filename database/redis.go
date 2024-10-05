package database

import (
	"api-dot/utils"
	"context"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     utils.GetConfig("REDIS_HOST") + ":" + utils.GetConfig("REDIS_PORT"),
		Password: utils.GetConfig("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return err
	}
	
	return nil
}

