package controllers

import (
	"GoTenancy/backend/config"
	"GoTenancy/backend/database/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type AdminController struct {
	Service services.UserService
	Ctx     iris.Context
	Session *sessions.Session
}

func (c *AdminController) Get() mvc.Result {
	view := mvc.View{
		Name: "admin/page/index.html",
		Data: iris.Map{
			"Title":   config.GetAppName(),
			"AppName": config.GetAppName(),
		},
	}
	return view
}
