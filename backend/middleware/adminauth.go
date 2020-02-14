package middleware

import (
	"GoTenancy/backend/session"
	"github.com/kataras/iris/v12"
)

func AdminAuth(ctx iris.Context) {
	if userId := session.Singleton().Start(ctx).GetInt64Default(session.UserIDKey, 0); userId <= 0 {
		ctx.Redirect("login")
		return
	}

	ctx.Next()
}
