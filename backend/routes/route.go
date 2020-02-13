package routes

import (
	"fmt"

	"GoTenancy/backend/controllers"
	"GoTenancy/backend/database"
	"GoTenancy/backend/middleware"
	"github.com/betacraft/yaag/irisyaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
)

func New(api *iris.Application) {
	app := api.Party("app.", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	{
		app.HandleDir("/static", "resources/app/static") // 加载租户端静态文件
		app.Get("/", func(ctx iris.Context) {            // 首页模块
			ctx.ViewLayout(iris.NoLayout) // 单页面应用无需模版  iris.NoLayout 或者 view.NoLayout
			if err := ctx.View("app/index.html"); err != nil {
				color.Red(fmt.Sprintf("商户端视图渲染出错 %v", err))
			}
		})

		v1 := app.Party("/v1")
		{
			v1.Use(irisyaag.New())
			v1.Post("/admin/login", controllers.UserLogin)
			v1.PartyFunc("/admin", func(app iris.Party) {
				casbinMiddleware := middleware.New(database.GetEnforcer())         //casbin for gorm                                                   // <- IMPORTANT, register the middleware.
				app.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
				app.Get("/logout", controllers.UserLogout).Name = "退出"

				app.PartyFunc("/users", func(users iris.Party) {
					users.Get("/", controllers.GetAllUsers).Name = "用户列表"
					users.Get("/{id:uint}", controllers.GetUser).Name = "用户详情"
					users.Post("/", controllers.CreateUser).Name = "创建用户"
					users.Put("/{id:uint}", controllers.UpdateUser).Name = "编辑用户"
					users.Delete("/{id:uint}", controllers.DeleteUser).Name = "删除用户"
					users.Get("/profile", controllers.GetProfile).Name = "个人信息"
				})
				app.PartyFunc("/roles", func(roles iris.Party) {
					roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
					roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
					roles.Post("/", controllers.CreateRole).Name = "创建角色"
					roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
					roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
				})
				app.PartyFunc("/permissions", func(permissions iris.Party) {
					permissions.Get("/", controllers.GetAllPermissions).Name = "权限列表"
					permissions.Get("/{id:uint}", controllers.GetPermission).Name = "权限详情"
					permissions.Post("/import", controllers.ImportPermission).Name = "导入权限"
					//permissions.Post("/", controllers.CreatePermission).Name = "创建权限"
					//permissions.Put("/{id:uint}", controllers.UpdatePermission).Name = "编辑权限"
					//permissions.Delete("/{id:uint}", controllers.DeletePermission).Name = "删除权限"
				})
			})
		}
	}

	// 管理端

	//admin := api.Party("admin.")
	//{

	//	//
	//	//api.OnAnyErrorCode(func(ctx iris.Context) {
	//	//	ctx.ViewData("Message", ctx.Values().
	//	//		GetStringDefault("message", "The page you're looking for doesn't exist"))
	//	//	ctx.View("admin/shared/error.html")
	//	//})
	//
	//	admin.HandleDir("/admin", "resources/admin")
	//	admin.Get("/", func(ctx iris.Context) { // 首页模块
	//		_ = ctx.View("admin/index.html")
	//	})
	//
	//	admin.Get("/apiDoc", func(ctx iris.Context) {
	//		_ = ctx.View("apiDoc/index.html")
	//	})
	//
	//	v1 := admin.Party("/v1")
	//	{
	//		v1.Use(irisyaag.New())
	//		v1.Post("/admin/login", controllers.AdminLogin)
	//		v1.Get("/admin/resetData", controllers.ResetData)
	//		v1.PartyFunc("/admin", func(admin iris.Party) {
	//			casbinMiddleware := middleware.New(database.GetEnforcer())           //casbin for gorm
	//			admin.Use(middleware.JwtHandler().Serve, casbinMiddleware.ServeHTTP) //登录验证
	//			admin.Get("/logout", controllers.UserLogout).Name = "退出"
	//
	//			admin.PartyFunc("/admin", func(users iris.Party) {
	//				users.Get("/profile", controllers.GetProfile).Name = "个人信息"
	//			})
	//			//admin.PartyFunc("/roles", func(roles iris.Party) {
	//			//	roles.Get("/", controllers.GetAllRoles).Name = "角色列表"
	//			//	roles.Get("/{id:uint}", controllers.GetRole).Name = "角色详情"
	//			//	roles.Post("/", controllers.CreateRole).Name = "创建角色"
	//			//	roles.Put("/{id:uint}", controllers.UpdateRole).Name = "编辑角色"
	//			//	roles.Delete("/{id:uint}", controllers.DeleteRole).Name = "删除角色"
	//			//})
	//		})
	//	}
	//
	//	api.Any("/payload", controllers.Payload)
	//}

	api.HandleDir("/doc", "resources/doc")
	api.Get("/", func(ctx iris.Context) { // 首页模块
		_ = ctx.View("doc/index.html")
	})

}
