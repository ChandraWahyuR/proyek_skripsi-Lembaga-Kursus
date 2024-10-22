package helper

import (
	"context"
	"fmt"
	"log"
	"skripsi/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	client *redis.Client
}

func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{client: client}
}

func (r *RedisHelper) BlacklistToken(ctx context.Context, tokenString string, expiresAt int64) error {
	return r.client.Set(ctx, tokenString, "blacklisted", time.Until(time.Unix(expiresAt, 0))).Err()
}

func (r *RedisHelper) IsTokenBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	val, err := r.client.Get(ctx, tokenString).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return val == "blacklisted", nil
}

func InitRedis(cfg *config.Config) *redis.Client {
	if cfg == nil || cfg.Redis.Host == "" || cfg.Redis.Port == 0 {
		// Log error jika konfigurasi Redis tidak ditemukan
		log.Fatal("Redis configuration is missing or incomplete")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Cek koneksi Redis
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	return rdb
}
