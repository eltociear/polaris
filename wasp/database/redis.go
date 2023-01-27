package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func newRedisClient() (RedisClient, error) {
	ctx := context.Background()

	redisdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisdb.Ping(ctx).Result()

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	redis := RedisClient{
		client: redisdb,
	}

	return redis, nil
}

func (r *RedisClient) Get(key string) ([]byte, error) {
	ctx := context.Background()
	val, err := r.get(ctx, key)
	return val, err
}

func (r *RedisClient) Set(key string, value []byte) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, time.Hour).Err()
}

func (r *RedisClient) Delete(key string) error {
	return nil
}

func (r *RedisClient) get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	switch {
	case err == redis.Nil:
		// key does not exist
		return nil, nil
	case err != nil:
		// Get failed
		return nil, err
	case len(val) == 0:
		// value is empty
		return nil, nil
	}
	return val, nil
}
