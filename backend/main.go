package main

import (
	"fmt"

	"GoTenancy/backend/config"
	"GoTenancy/backend/database"
	"GoTenancy/backend/database/models"
	"GoTenancy/backend/logs"
	"GoTenancy/backend/routes"
	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func NewApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel(config.GetAppLoggerLevel())

	app.RegisterView(
		iris.HTML("resources", ".html"). // 加载模版文件
							Layout("shared/layout.html"). // 增加布局模版
							Reload(true),                 // 增加静态文件重载
	)
	admin := app.Party("admin.")
	admin.HandleDir("/admin", "resources/admin") // 注册管理端静态文件
	mvc.Configure(admin, routes.AdminMVC)        // 注册管理端 mvc

	db := database.GetGdb()
	db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware. //api 文档配置
		On:       true,
		DocTitle: config.GetAppName(),
		DocPath:  "./resources/apiDoc/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": config.GetAppUrl(),
			"Staging":    "",
		},
	})

	routes.New(app) //注册 app 路由

	return app
}

func main() {
	f := logs.NewLog()
	defer f.Close()

	app := NewApp()
	//app.Logger().SetOutput(f) //记录日志
	err := app.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf()))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
