package apiserver

import (
	"net/http"

	"minibox.ai/pkg/sessions"
	"minibox.ai/pkg/utils/tmpl"
)

func homeMiddleware(next http.Handler) http.Handler {
	website := mountWebsiteHandler("../website/build/minibox")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			session := sessions.GetSession(r)
			cur_user := getCurrentUser(session)
			if cur_user == nil { // 没有用户信息
				// 显示一个空白
				website.ServeHTTP(w, r)
				// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
				return
			}

			// 显示用户主页
			tmpl.RenderTemplate(w, "home.html.tmpl", cur_user)
		} else {
			// 继续下一层处理，如 Namespace 分析
			next.ServeHTTP(w, r)
		}
	})
}
