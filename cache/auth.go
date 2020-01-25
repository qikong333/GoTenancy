package cache

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// Auth 用于获取/设置身份验证相关密钥
type Auth struct{}

// Exists 从缓存返回身份验证。
func (x *Auth) Exists(key string, v interface{}) error {
	s, err := rc.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	dec := gob.NewDecoder(strings.NewReader(s))
	return dec.Decode(v)
}

// Set 缓存此密钥（身份验证）30 分钟。
func (x *Auth) Set(key string, v interface{}, expiration time.Duration) error {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return err
	}

	return rc.Set(key, buf.String(), expiration).Err()
}
