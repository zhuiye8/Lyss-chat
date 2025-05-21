package db

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/config"
)

// Redis 表示 Redis 连接
type Redis struct {
	Client *redis.Client
}

// NewRedis 创建一个新的 Redis 连接
func NewRedis(cfg config.RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

// Close 关闭 Redis 连接
func (r *Redis) Close() error {
	return r.Client.Close()
}
