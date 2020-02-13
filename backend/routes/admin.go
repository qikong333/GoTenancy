package routes

import (
	"time"

	"GoTenancy/backend/config"
	"GoTenancy/backend/controllers"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

func AdminMVC(app *mvc.Application) {
	app.Router.Use(
		func(ctx iris.Context) {
			ctx.Application().Logger().Infof("Path: %s Method: %s", ctx.Path(), ctx.Method())
			ctx.Next()
		},
	)

	app.Handle(new(controllers.AdminAuthController))

	party := app.Party("/home")
	party.Register(
		sessions.New(sessions.Config{
			Cookie:  config.GetAppCookieNameForSessionID(),
			Expires: 24 * time.Hour,
		}),
	)
	//party.Router.Use(middleware.AdminAuth)
	party.Handle(new(controllers.AdminController))
}
