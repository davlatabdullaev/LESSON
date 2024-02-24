package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"test/config"
	"test/storage"
	"time"
)

type Store struct {
	db *redis.Client
}

func New(cfg config.Config) storage.IRedisStorage {
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})

	return Store{
		db: redisClient,
	}
}

func (s Store) SetX(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	statusCmd := s.db.SetEx(ctx, key, value, duration)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}

	return nil
}

func (s Store) Get(ctx context.Context, key string) interface{} {
	return s.db.Get(ctx, key).Val()
}