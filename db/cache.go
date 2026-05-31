package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"singo/util"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	client := redis.NewClient(&redis.Options{
		Addr:       viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password:   viper.GetString("redis.password"),
		DB:         viper.GetInt("redis.database"),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		util.Log().Panic("连接Redis不成功", err)
	}

	RedisClient = client
}
