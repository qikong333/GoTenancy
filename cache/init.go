package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/snowlyg/GoTenancy/queue"
)

var rc *redis.Client

func init() {
	host := os.Getenv("REDIS_ADDR")
	if len(host) == 0 {
		host = "127.0.0.1:6379"
	}

	c := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: os.Getenv("REDIS_KEY"),
		DB:       0, // use default DB
	})

	if _, err := c.Ping().Result(); err != nil {
		log.Fatal("unable to connect to redis", err)
	}

	rc = c
}

// New 通过 queue.New 函数初始化队列服务。
// queueProcessor 标识表明是否实例执行发布/订阅 订阅者。只能有一个订阅者。
// ex 参数 map[queue.TaskID]queue.Executor 允许你给自定义的任务提供
// 自定义执行者。 TaskExecutor 必须满足
// 此接口。
// 	type TaskExecutor interface {
// 		Run(t QueueTask) error
// 	}
func New(queueProcessor, isDev bool, ex map[queue.TaskID]queue.TaskExecutor) {
	queue.New(rc, isDev, ex)

	if queueProcessor {
		go queue.SetAsSubscriber()
	}
}
