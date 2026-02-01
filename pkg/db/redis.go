package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"test/config"
	"test/pkg/logger"
	"time"
)

type RedisTemplate struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisTemplates() (*RedisTemplate, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GConfig.Redis.Addr,
		Password: config.GConfig.Redis.Passward,
		DB:       config.GConfig.Redis.DB,
		PoolSize: config.GConfig.Redis.PoolSize,
		//超时配置
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolTimeout:  2 * time.Second,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis连接错误,%s", err)
	}
	logger.Info("redis连接成功...")
	return &RedisTemplate{client: client,
		ctx: context.Background()}, nil
}
func (r *RedisTemplate) SetJson(key string, o interface{}, expire time.Duration) error {
	marshal, err := json.Marshal(o)
	if err != nil {
		logger.Error("redis插入序列化失败!", zap.String("key", key))
		return nil
	}
	err2 := r.client.Set(r.ctx, key, marshal, 0).Err()
	if err2 != nil {
		logger.Error("redis插入失败!", zap.String("key", key))
		return err2
	}
	return nil
}
func (r *RedisTemplate) GetJson(key string) (interface{}, error) {
	bytes, err := r.client.Get(r.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("redis can not find the key:%s", key)
		}
		return nil, fmt.Errorf("redis get %s failed", key)
	}
	var o interface{}
	err = json.Unmarshal(bytes, &o)
	if err != nil {
		fmt.Println("反序列化失败")
		return nil, fmt.Errorf("获取结果反序列化失败%s", key)
	}
	//通过泛型适配
	return o, nil
}
