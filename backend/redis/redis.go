package redis

import (
	"sync"
	"time"

	"GoTenancy/backend/config"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

var db *redis.Database
var once sync.Once

/**
*设置数据库连接
*@param diver string
 */
func Singleton() *redis.Database {
	once.Do(func() {
		db = redis.New(redis.Config{
			Network:   "tcp",
			Addr:      config.GetRedisAddr(),
			Timeout:   time.Duration(30) * time.Second,
			MaxActive: 10,
			Password:  config.GetRedisPwd(),
			Database:  "",
			Prefix:    "",
			Delim:     "-",
			Driver:    redis.Redigo(), // redis.Radix() can be used instead.
		})
	})

	return db
}
