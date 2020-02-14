package validates

type CreateUpdateAdminRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
	Name     string `json:"name" validate:"required,gte=6,lte=50"  comment:"名称"`
	RoleIds  []uint `json:"role_ids"  validate:"required" comment:"角色"`
}

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required,gte=5,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required"  comment:"密码"`
}

// 登陆验证表单
func (alr *AdminLoginRequest) Valid() string {
	return BaseValid(alr)
}

// 创建更新验证表单
func (alr *CreateUpdateAdminRequest) Valid() string {
	return BaseValid(alr)
}
