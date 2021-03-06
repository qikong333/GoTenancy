package controllers

import (
	"GoTenancy/backend/database/services"
	"github.com/kataras/iris/v12"
)

type Response struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

type Lists struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func ApiResource(status bool, objects interface{}, msg string) (r *Response) {
	r = &Response{Status: status, Data: objects, Msg: msg}
	return
}

type BaseAdminController struct {
	Service services.UserService
	Ctx     iris.Context
}
