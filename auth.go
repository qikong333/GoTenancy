package GoTenancy

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/snowlyg/GoTenancy/cache"
	"github.com/snowlyg/GoTenancy/data"
	"github.com/snowlyg/GoTenancy/model"
)

// Auth 表示一个已认证用户
type Auth struct {
	AccountID int64
	UserID    int64
	Email     string
	Role      model.Roles
}

// Authenticator 中间件应用于认证请求
//
// 有 4 种方法去认证一个请求:
// 1. 通过一个 X-API-KEY HTTP 请求头。
// 2. 通过一个 "key=token" 请求参数。
// 3. 通过一个 X-API-KEY cookie。
// 4. 通过基础认证
//
// 对于将最小角色设置为 model.RolePublic 的路由，将不会执行身份验证。
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mr := ctx.Value(ContextMinimumRole).(model.Roles)

		key, pat, err := extractKeyFromRequest(r)
		// 没有认证，或者是一个错误
		if mr > model.RolePublic {
			if len(key) == 0 || err != nil {
				http.Redirect(w, r, "/users/login", http.StatusSeeOther)
				return
			}
		}

		ca := &cache.Auth{}

		// key 已经被缓存
		var a Auth
		if err := ca.Exists(key, &a); err != nil {
			log.Println("error while trying to get cache auth", err)
		}

		if len(a.Email) > 0 {
			ctx = context.WithValue(ctx, ContextAuth, a)
		} else {
			// 如果是公共路由，则不需要任何验证
			if mr == model.RolePublic {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			db, ok := ctx.Value(ContextDatabase).(*data.DB)
			if !ok {
				http.Error(w, "database not available", http.StatusUnauthorized)
				return
			}

			id, t := model.ParseToken(key)
			acct, usr, err := db.Users.Auth(id, t, pat)
			if err != nil {
				er := fmt.Sprintf("invalid token key: %v", err)
				http.Error(w, er, http.StatusUnauthorized)
				return
			}

			a.AccountID = acct.ID
			a.Email = usr.Email
			a.UserID = usr.ID
			a.Role = usr.Role

			// 保存到缓存
			ca.Set(key, a, 30*time.Minute)

			ctx = context.WithValue(ctx, ContextAuth, a)
		}

		// 认证请求
		if a.Role < mr {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractKeyFromRequest(r *http.Request) (key string, pat bool, err error) {
	// 首先，让我们看看是否在 HTTP 标头中存在 X-API-KEY
	key = r.Header.Get("X-API-KEY")
	if len(key) > 0 {
		return
	}

	// 检查查询字符串
	key = r.URL.Query().Get("key")
	if len(key) > 0 {
		return
	}

	// 检查 cookie
	ck, er := r.Cookie("X-API-KEY")
	if er != nil {
		//如果是 ErrNoCookie，我们必须继续，
		//否则这是一个合法的错误
		if er != http.ErrNoCookie {
			err = er
			return
		}
	} else {
		key = ck.Value
		return
	}

	// 检查是否支持基础认证
	authorization := r.Header.Get("Authorization")
	s := strings.SplitN(authorization, " ", 2)
	if len(s) != 2 {
		err = fmt.Errorf("invalid basic authentication format: %s - you must provide Basic base64token", authorization)
		return
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		err = fmt.Errorf("invalid basic authentication format: %s - you must provide Basic base64token", authorization)
		return
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		err = fmt.Errorf("invalid basic authentication, your token should be _:access_token - got %s", string(b))
		return
	}

	key = pair[1]
	pat = true

	return
}
