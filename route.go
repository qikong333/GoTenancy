package GoTenancy

import (
	"net/http"

	"github.com/snowlyg/GoTenancy/model"
)

// Route 结构体表示具有可选中间件的 Web 处理程序。
type Route struct {
	// 中间件
	WithDB           bool
	Logger           bool
	EnforceRateLimit bool
	AllowCrossOrigin bool

	// 认证
	MinimumRole model.Roles

	Handler http.Handler
}

// NewError 方法返回一个拥有错误和状态码 Respond 的 Route 对象。
func NewError(err error, statusCode int) *Route {
	return &Route{
		Logger:      true,
		MinimumRole: model.RolePublic,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Respond(w, r, statusCode, err)
		}),
	}
}

func notFound(w http.ResponseWriter) {
	http.Error(w, "not found", http.StatusNotFound)
}
