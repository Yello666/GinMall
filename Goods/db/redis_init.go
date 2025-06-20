package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var Rdb *redis.Client

func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.64.2:6379", // redis地址
		Password: "123456",            // 没有密码
		DB:       0,                   // 默认DB
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis连接失败:", err)
	}
	return err
}
