package component

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/mingslife/bone"

	"elf-server/pkg/conf"
)

type Cache struct {
	Logger *DatabaseLogger `inject:""`
	RDB    *redis.Client
}

func (*Cache) Name() string {
	return "component.cache"
}

func (*Cache) Init() error {
	return nil
}

func (c *Cache) Register() error {
	cfg := conf.GetConfig()
	c.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPwd,
		DB:       cfg.RedisDb,
	})
	err := c.RDB.Ping(context.Background()).Err()
	return err
}

func (c *Cache) Unregister() error {
	return nil
}

var _ bone.Component = (*Cache)(nil)
