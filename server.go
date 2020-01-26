package GoTenancy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/internal/config"
)

func init() {
	if err := config.LoadFromFile(); err != nil {
		log.Println(err)
	}

	if len(config.Current.StripeKey) > 0 {
		SetStripeKey(config.Current.StripeKey)
	}

	if len(config.Current.Plans) > 0 {
		for _, p := range config.Current.Plans {
			if p.Params == nil {
				p.Params = make(map[string]interface{})
			}
			data.AddPlan(p)
		}
	}
}

// Server 结构体是后端启动的端点
// 负责将路由请求转到处理程序
type Server struct {
	DB              data.DB
	Logger          func(http.Handler) http.Handler
	Authenticator   func(http.Handler) http.Handler
	Throttler       func(http.Handler) http.Handler
	RateLimiter     func(http.Handler) http.Handler
	Cors            func(http.Handler) http.Handler
	StaticDirectory string
	Routes          map[string]*Route
}

// NewServer 方法返回一个拥有所有可用中间件的项目服务。
// 只有顶级的路由需要作为参数传递。
// 三个内置实现如下：
// 1. users: 用户管理（登陆, 注册, 认证, 详情等等）。
// 2. billing: 用于全功能计费流程（从免费转换为付费、更改计划、获取发票等）。
// 3. webhooks: 允许用户订阅事件（你可以通过 GoTenancy.SendWebhook 取消 webhook 事件）。
// 你可以在项目中，简单的重写默认的内置实现：
// 	routes := make(map[string]*GoTenancy.Route)
// 	routes["billing"] = &GoTenancy.Route{Handler: billing.Route}

func NewServer(routes map[string]*Route) *Server {
	// 如果 users, billing 和 webhooks 路由没有被定义，
	// GoTenancy 将会默认实现它们。
	if _, ok := routes["users"]; !ok {
		routes["users"] = newUser()
	}

	if _, ok := routes["billing"]; !ok {
		routes["billing"] = newBilling()
	}

	if _, ok := routes["webhooks"]; !ok {
		routes["webhooks"] = newWebhook()
	}

	return &Server{
		Logger:          Logger,
		Authenticator:   Authenticator,
		Throttler:       Throttler,
		RateLimiter:     RateLimiter,
		Cors:            Cors,
		StaticDirectory: "/public/",
		Routes:          routes,
	}
}

// 如果没有路由，将会返回一个错误
// 静态文件默认保存在 "/public/" 目录中。你可以设置 StaticDirectory 属性来改变默认的静态文件目录：
// 	mux := GoTenancy.NewServer(routes)
// 	mux.StaticDirectory = "/files/"
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if strings.HasPrefix(r.URL.Path, s.StaticDirectory) {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextOriginalPath, r.URL.Path)

	isJSON := strings.ToLower(r.Header.Get("Content-Type")) == "application/json"
	ctx = context.WithValue(ctx, ContextContentIsJSON, isJSON)

	var next *Route
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	if r, ok := s.Routes[head]; ok {
		next = r
	} else if catchall, ok := s.Routes["__catchall__"]; ok {
		next = catchall
	} else {
		next = NewError(fmt.Errorf("path not found"), http.StatusNotFound)
	}

	if next.WithDB {
		ctx = context.WithValue(ctx, ContextDatabase, s.DB)
	}

	ctx = context.WithValue(ctx, ContextMinimumRole, next.MinimumRole)

	// 验证所有处理程序
	next.Handler = s.Authenticator(next.Handler)

	if next.Logger {
		next.Handler = s.Logger(next.Handler)
	}

	if next.EnforceRateLimit {
		next.Handler = s.RateLimiter(next.Handler)
		next.Handler = s.Throttler(next.Handler)
	}

	// 是否允许此路由的跨域请求
	if next.AllowCrossOrigin {
		next.Handler = s.Cors(next.Handler)
	}

	next.Handler.ServeHTTP(w, r.WithContext(ctx))
}
