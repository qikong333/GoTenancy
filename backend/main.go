package main

import (
	"fmt"

	"GoTenancy/backend/config"
	"GoTenancy/backend/database"
	"GoTenancy/backend/database/models"
	"GoTenancy/backend/libs"
	"GoTenancy/backend/logs"
	"GoTenancy/backend/redis"
	"GoTenancy/backend/routes"
	"GoTenancy/backend/tasks"
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
	models.AutoMigrate()

	iris.RegisterOnInterrupt(func() {
		redisErr := redis.Singleton().Close()
		if redisErr != nil {
			color.Red(fmt.Sprintf("redis: %v", redisErr))
		}
		dbErr := db.Close()
		if dbErr != nil {
			color.Red(fmt.Sprintf("db: %v", dbErr))
		}
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
	tasks.New()

	return app
}

func main() {
	f := logs.NewLog()
	defer f.Close()

	defer libs.StopTask()
	defer redis.Singleton().Close()

	app := NewApp()
	//app.Logger().SetOutput(f) //记录日志
	err := app.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf()))
	if err != nil {
		color.Red(fmt.Sprintf("项目运行结束: %v", err))
	}
}
