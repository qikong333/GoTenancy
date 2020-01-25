package GoTenancy

import (
	"path"
	"strings"
)

// ShiftPath 方法分割请求 URL 的头部和尾部
// 这对于在 ServeHTTP 函数内部执行路由很有用。
// 	package yourapp
// 	import (
// 		"net/http"
// 		"github.com/dstpierre/GoTenancy"
// 	)
//
// 	func main() {
// 		routes := make(map[string]*GoTenancy.Route)
// 		routes["speak"] = &GoTenancy.Route{Handler: speak{}}
// 		mux := GoTenancy.NewServer(routes)
// 		http.ListenAndServe(":8080", mux)
// 	}
//
// 	type speak struct{}
// 	func (s speak) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 		var head string
// 		head, r.URL.Path = GoTenancy.ShiftPath(r.URL.Path)
// 		if head == "loud" {
// 			s.scream(w, r)
// 		} else {
// 			s.whisper(w, r)
// 		}
// 	}
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
