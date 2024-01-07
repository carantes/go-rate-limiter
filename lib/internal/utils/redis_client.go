package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)

	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisClient) Set(key string, obj interface{}, ttl time.Duration) error {
	value, _ := json.Marshal(obj)

	return r.client.Set(r.ctx, key, string(value), ttl).Err()
}

func (r *RedisClient) Get(key string, obj interface{}) error {
	data, err := r.client.Get(r.ctx, key).Result()

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), obj); err != nil {
		return err
	}

	return nil
}
