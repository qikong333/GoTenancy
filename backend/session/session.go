package session

import (
	"sync"
	"time"

	"GoTenancy/backend/config"
	"GoTenancy/backend/redis"
	"github.com/kataras/iris/v12/sessions"
)

const UserIDKey = "UserID"

var (
	sess *sessions.Sessions
	once sync.Once
)

func Singleton() *sessions.Sessions {
	once.Do(func() {
		sess = sessions.New(
			sessions.Config{
				Cookie:       config.GetAppCookieNameForSessionID(),
				AllowReclaim: true,
				Expires:      4 * time.Hour,
			},
		)
		sess.UseDatabase(redis.Singleton())
	})
	return sess
}
