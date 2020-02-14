package controllers

import (
	"GoTenancy/backend/config"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminController struct {
	BaseAdminController
}

func (c *AdminController) Get() mvc.Result {
	c.Ctx.ViewLayout(iris.NoLayout)
	view := mvc.View{
		Name: "admin/index.html",
		Data: iris.Map{
			"Title":   config.GetAppName(),
			"AppName": config.GetAppName(),
		},
	}
	return view
}
