package session

import (
	"sync"

	"GoTenancy/backend/config"
	"github.com/kataras/iris/v12/sessions"
)

var (
	sess *sessions.Sessions
	once sync.Once
)

func Single() *sessions.Sessions {
	once.Do(func() {
		sess = sessions.New(sessions.Config{Cookie: config.GetAppCookieNameForSessionID(), AllowReclaim: true})
	})
	return sess
}
