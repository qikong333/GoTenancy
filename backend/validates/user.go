package validates

type CreateUpdateUserRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
	Name     string `json:"name" validate:"required,gte=6,lte=50"  comment:"名称"`
	RoleIds  []uint `json:"role_ids"  validate:"required" comment:"角色"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

// 登陆表单验证
func (alr *LoginRequest) Valid() string {
	return BaseValid(alr)
}

// 新建更新表单验证
func (alr *CreateUpdateUserRequest) Valid() string {
	return BaseValid(alr)
}
