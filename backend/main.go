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
)

func NewApp() *iris.Application {
	api := iris.New()
	api.Logger().SetLevel(config.GetAppLoggerLevel())
	api.RegisterView(iris.HTML("resources", ".html"))

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

	routes.App(api) //注册 app 路由

	return api
}

func main() {
	f := logs.NewLog()
	defer f.Close()

	api := NewApp()
	api.Logger().SetOutput(f) //记录日志
	err := api.Run(iris.Addr(config.GetAppUrl()), iris.WithConfiguration(config.GetIrisConf()))
	if err != nil {
		color.Yellow(fmt.Sprintf("项目运行结束: %v", err))
	}
}
