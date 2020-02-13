package routes

import (
	"fmt"

	"GoTenancy/backend/controllers"
	"GoTenancy/backend/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type prefixedLogger struct {
	prefix string
}

func (s *prefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.prefix, msg)
}

func AdminMVC(app *mvc.Application) {
	// You can use normal middlewares at MVC apps of course.
	app.Router.Use(
		func(ctx iris.Context) {
			ctx.Application().Logger().Infof("Path: %s", ctx.Path())
			ctx.Next()
		},
		middleware.BasicAuth,
	)

	// Register dependencies which will be binding to the controller(s),
	// can be either a function which accepts an iris.Context and returns a single value (dynamic binding)
	// or a static struct value (service).
	app.Register(
		sessions.New(sessions.Config{}).Start,
		&prefixedLogger{prefix: "DEV"},
	)

	// GET: http://localhost:8080/basic
	// GET: http://localhost:8080/basic/custom
	// GET: http://localhost:8080/basic/custom2
	app.Handle(new(controllers.AdminController))

	// All dependencies of the parent *mvc.Application
	// are cloned to this new child,
	// thefore it has access to the same session as well.
	// GET: http://localhost:8080/basic/sub
	//app.Party("/sub").
	//	Handle(new(basicSubController))
}
