package validates

type TenancyRequest struct {
	Name string `json:"name" validate:"required,gte=6,lte=50"  comment:"名称"`
}


// 新建更新表单验证
func (alr *TenancyRequest) Valid() string {
	return BaseValid(alr)
}
