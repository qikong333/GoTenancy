package middleware

import (
	"GoTenancy/backend/config"
	"github.com/kataras/iris/v12/middleware/basicauth"
)

// BasicAuth middleware sample.
var BasicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		config.GetAdminUserName(): config.GetAdminPwd(),
	},
})
