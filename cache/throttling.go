package cache

import (
	"fmt"
	"time"
)

// Throttle 增加特定密钥的请求计数，并设置过期（如果是新期间）。
func Throttle(key string, expire time.Duration) (int64, error) {
	key = fmt.Sprintf("%s_t", key)
	return increaseThrottle(key, expire)
}

// RateLimit 增加特定密钥的请求计数，并设置过期（如果是新期间）。
func RateLimit(key string, expire time.Duration) (int64, error) {
	key = fmt.Sprintf("%s_rl", key)
	return increaseThrottle(key, expire)
}

func increaseThrottle(key string, expire time.Duration) (int64, error) {
	i, err := rc.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	if i == 1 {
		// 密钥已创建，我们设置过期
		ok, err := rc.Expire(key, expire).Result()
		if err != nil {
			// 尝试删除密钥
			if _, e := rc.Del(key).Result(); err != nil {
				return 0, fmt.Errorf("unable to remove key %s: %s and expire failed: %s", key, e.Error(), err.Error())
			}
			return 0, err
		} else if !ok {
			return 0, fmt.Errorf("unable to set expiration on key %s", key)
		}
	}

	return i, nil
}

// GetThrottleExpiration 返回密钥过期之前的限制持续时间。
func GetThrottleExpiration(key string) (time.Duration, error) {
	key = fmt.Sprintf("%s_t", key)

	return rc.TTL(key).Result()
}

// GetRateLimitExpiration 返回密钥过期之前的限制持续时间。
func GetRateLimitExpiration(key string) (time.Duration, error) {
	key = fmt.Sprintf("%s_rl", key)

	return rc.TTL(key).Result()
}
