package routes

import (
	"GoTenancy/backend/controllers"
	"GoTenancy/backend/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
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
	party.Router.Use(middleware.AdminAuth)
	party.Handle(new(controllers.AdminController))
}
