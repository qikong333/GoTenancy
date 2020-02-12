package validates

type AppRequest struct {
	Name string `json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
}

// 新建更新表单验证
func (alr *AppRequest) Valid() string {
	return BaseValid(alr)
}
