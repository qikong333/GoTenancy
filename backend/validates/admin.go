package validates

import (
	"github.com/go-playground/validator/v10"
)

type CreateUpdateAdminRequest struct {
	UserName string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
	RoleIds  []uint `json:"role_ids"  validate:"required" comment:"角色"`
}

type AdminLoginRequest struct {
	UserName string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

// 验证表单
func (alr *AdminLoginRequest) Valid() string {
	var formErrs string
	err := Validate.Struct(*alr)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(ValidateTrans) {
			if len(e) > 0 {
				formErrs += formErrs + ";"
			}
		}
	}
	return formErrs
}
