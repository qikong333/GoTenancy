package GoTenancy

import (
	"context"
	"net/http"
)

// Language 是一个处理名为 "lng".Language 语言 cookie 的中间件
// 它在 HTML 模版和Go 代码中使用 Translate 方法的时候使用.
// 你需要在一个名为  languagepack 的文件夹中创建一个语言文件 (i.e. en.json, fr.json)。
func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lng := "en"
		ck, err := r.Cookie("lng")
		if err == nil {
			lng = ck.Value
		}

		ctx := context.WithValue(r.Context(), ContextLanguage, lng)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
