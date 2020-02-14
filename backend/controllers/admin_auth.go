package controllers

import (
	"fmt"

	"GoTenancy/backend/config"
	"GoTenancy/backend/database/models"
	"GoTenancy/backend/routepath"
	"GoTenancy/backend/session"
	"GoTenancy/backend/validates"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type AdminAuthController struct {
	BaseAdminController
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

func (c *AdminAuthController) PostLogin() {

	aul := new(validates.AdminLoginRequest)
	if err := c.Ctx.ReadJSON(aul); err != nil {
		color.Red(fmt.Sprintf("ReadJSON:%v", err))
		_, _ = c.Ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	if formErrs := aul.Valid(); len(formErrs) > 0 {
		color.Red(fmt.Sprintf("Valid:%v", formErrs))
		_, _ = c.Ctx.JSON(ApiResource(false, nil, formErrs))
		return
	}

	admin := models.NewAdmin(0, aul.Username)
	admin.GetAdminByUserName()

	status, msg := admin.CheckLogin(aul.Password)
	if status {
		c.Ctx.Application().Logger().Infof("%s 登录系统", aul.Username)
		session.Singleton().Start(c.Ctx).Set(session.UserIDKey, int64(admin.ID))
	}
	_, _ = c.Ctx.JSON(ApiResource(status, nil, msg))
	return

}

func (c *AdminAuthController) AnyLogout() {
	if b := session.Singleton().Start(c.Ctx).Delete(session.UserIDKey); b {
		_, _ = c.Ctx.JSON(ApiResource(true, nil, "退出登陆"))
		return
	}
	_, _ = c.Ctx.JSON(ApiResource(false, nil, "退出失败"))
}

func (c *AdminAuthController) GetResert() {
	models.DelAllData()
	routes := routepath.GetRoutes(c.Ctx.Application().GetRoutesReadOnly())
	models.CreateSystemData(routes) // 初始化系统数据 账号，角色，权限
	c.Ctx.StatusCode(iris.StatusOK)
	_, _ = c.Ctx.JSON(ApiResource(true, routes, "重置数据成功"))
}
