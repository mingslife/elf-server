package components

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var Cache *cacheComponent

type cacheComponent struct {
	ctx context.Context
	rdb *redis.Client
}

func (c *cacheComponent) Get(key string) string {
	return c.rdb.Get(c.ctx, key).Val()
}

func (c *cacheComponent) Set(key, value string, expiration time.Duration) {
	c.rdb.Set(c.ctx, key, value, expiration)
}

func (c *cacheComponent) Del(key string) {
	c.rdb.Del(c.ctx, key)
}

func InitCache(host string, port int, password string, db int) {
	c := &cacheComponent{}

	c.ctx = context.Background()
	c.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})
	if err := c.rdb.Ping(c.ctx).Err(); err != nil {
		panic(err)
	}

	Cache = c
}

type CacheCaptchaStore struct {
	expiration time.Duration
}

func (s *CacheCaptchaStore) key(id string) string {
	return fmt.Sprintf("captcha:%s", id)
}

func (s *CacheCaptchaStore) Set(id string, digits []byte) {
	key := s.key(id)
	ns := make([]byte, len(digits))
	for i, d := range digits {
		ns[i] = '0' + d
	}
	value := string(ns)
	Cache.Set(key, string(value), s.expiration)
}

func (s *CacheCaptchaStore) Get(id string, clear bool) (digits []byte) {
	key := s.key(id)
	value := Cache.Get(key)
	ns := []rune(value)
	digits = make([]byte, len(ns))
	for i, n := range ns {
		digits[i] = byte(n - '0')
	}
	if clear {
		Cache.Del(key)
	}
	return
}

func NewCacheCaptchaStore(expiration time.Duration) *CacheCaptchaStore {
	s := &CacheCaptchaStore{}
	s.expiration = expiration
	return s
}
