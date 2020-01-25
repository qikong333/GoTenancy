package GoTenancy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/snowlyg/GoTenancy/cache"
)

// RateLimiter 是一个限制短时间内发起太多请求的中间件。
// 如果达到每个用户允许最大的请求次数， 将会返回一个 StatusTooManyRequests 错误。
// 简单来说，如果用户在同一个 "Retry-After" HTTP 请求头下，达到了最大的请求次数，
// 那么用户需要1分钟后才可以发送一个新的请求。
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var keys Auth

		ctx := r.Context()
		v := ctx.Value(ContextAuth)
		if v == nil {
			keys = Auth{}
		} else {
			a, ok := v.(Auth)
			if ok {
				keys = a
			}
		}

		key := fmt.Sprintf("%v", keys.AccountID)

		// TODO: Make this configurable
		count, err := cache.RateLimit(key, 1*time.Minute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Make this configurable
		if count >= 60 {
			d, err := cache.GetThrottleExpiration(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if d.Seconds() > 0 {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(d.Seconds())))
			}
			http.Error(w, fmt.Sprintf("you've reached your rate limit, retry in %v", d), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
