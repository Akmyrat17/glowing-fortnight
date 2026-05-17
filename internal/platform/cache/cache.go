package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/boilerplate/internal/config"
)

func New(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	return client, nil
}
