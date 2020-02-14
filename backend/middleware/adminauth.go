package middleware

import (
	"fmt"

	"GoTenancy/backend/session"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
)

func AdminAuth(ctx iris.Context) {
	if userId := session.Singleton().Start(ctx).GetInt64Default(session.UserIDKey, 0); userId <= 0 {
		color.Red(fmt.Sprintf("UserIDKey %d", userId))
		ctx.Redirect("login")
		return
	}

	ctx.Next()
}
