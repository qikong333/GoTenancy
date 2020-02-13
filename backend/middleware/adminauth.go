package middleware

import (
	session2 "GoTenancy/backend/session"
	"github.com/kataras/iris/v12"
)

func AdminAuth(ctx iris.Context) {
	// Check if user is authenticated
	if ctx.Path() == "/login" && ctx.Method() == "GET" {
		ctx.Next()
	}

	if auth, _ := session2.Single().Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.Redirect("login")
		return
	}

	ctx.Next()
}
