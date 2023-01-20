package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type Redis struct {
	client *redis.Client
}

func New(redisHost string, redisPort int, password string, db int) Redis {
	ctx := context.Background()

	redisUri := fmt.Sprintf("%s:%d", redisHost, redisPort)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: "",
		DB:       0, //use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	redis := Redis{client: redisClient}

	return redis
}

func NewFromEnv() Redis {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(ctx).Result()

	if err != nil {
		panic(err)
	}

	redis := Redis{client: redisClient}

	return redis
}

func (r *Redis) Get(key string) (interface{}, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}
