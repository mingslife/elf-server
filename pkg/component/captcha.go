package component

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dchest/captcha"
	"github.com/mingslife/bone"
)

type Captcha struct {
	Logger     *DatabaseLogger `inject:""`
	Cache      *Cache          `inject:"component.cache"`
	expiration time.Duration
}

func (*Captcha) Name() string {
	return "component.captcha"
}

func (*Captcha) Init() error {
	return nil
}

func (c *Captcha) Register() error {
	if c.Cache == nil {
		return errors.New("cache is nil")
	}
	c.expiration = 10 * time.Minute
	captcha.SetCustomStore(c)
	return nil
}

func (c *Captcha) Unregister() error {
	return nil
}

func (*Captcha) key(id string) string {
	return fmt.Sprintf("captcha:%s", id)
}

func (c *Captcha) Set(id string, digits []byte) {
	key := c.key(id)
	ns := make([]byte, len(digits))
	for i, d := range digits {
		ns[i] = '0' + d
	}
	value := string(ns)
	c.Cache.RDB.Set(context.Background(), key, string(value), c.expiration)
}

func (c *Captcha) Get(id string, clear bool) (digits []byte) {
	key := c.key(id)
	value := c.Cache.RDB.Get(context.Background(), key).Val()
	ns := []rune(value)
	digits = make([]byte, len(ns))
	for i, n := range ns {
		digits[i] = byte(n - '0')
	}
	if clear {
		c.Cache.RDB.Del(context.Background(), key)
	}
	return
}

var _ bone.Component = (*Captcha)(nil)
