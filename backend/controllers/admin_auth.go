package controllers

import (
	"GoTenancy/backend/config"
	"GoTenancy/backend/database/models"
	"GoTenancy/backend/database/services"
	"GoTenancy/backend/validates"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

const userIDKey = "UserID"

type AdminAuthController struct {
	Service services.UserService
	Ctx     iris.Context
	Session *sessions.Session
}

func (c *AdminAuthController) Get() {
	c.Ctx.Redirect("home")
}

func (c *AdminAuthController) GetLogin() mvc.Result {
	c.Ctx.ViewLayout(iris.NoLayout)
	view := mvc.View{
		Name: "admin/page/login/index.html",
		Data: iris.Map{
			"Title":   config.GetAppName(),
			"AppName": config.GetAppName(),
		},
	}
	return view
}

func (c *AdminController) PostLogin() {
	aul := new(validates.AdminLoginRequest)
	if err := c.Ctx.ReadJSON(aul); err != nil {
		c.Ctx.Redirect("login")
	}

	if formErrs := aul.Valid(); len(formErrs) > 0 {
		c.Ctx.Redirect("login")
	}

	admin := models.NewAdmin(0, aul.UserName)
	admin.GetAdminByUserName()

	status, _ := admin.CheckLogin(aul.Password)
	if status {
		c.Ctx.Application().Logger().Infof("%s 登录系统", aul.UserName)
		c.Session.Set(userIDKey, admin.ID)
	} else {
		c.Ctx.Redirect("login")
	}

}
