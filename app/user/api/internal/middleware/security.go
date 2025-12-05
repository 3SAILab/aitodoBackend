package middleware

import "net/http"

// SecurityHeadersMiddleware 为所有响应统一添加安全相关响应头
func SecurityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next(w, r)
	}
}

// CSRFMiddleware 实现双重提交 Cookie 的 CSRF 校验：
// 对于有副作用的请求（POST/PUT/PATCH/DELETE），要求：
//   Cookie 中存在 csrfToken，且 Header "X-CSRF-Token" 与之相同。
func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			// 登录 / 刷新 / 登出接口依赖 Cookie，本身不做 CSRF 校验
			switch r.URL.Path {
			case "/users/login", "/users/refresh-token", "/users/logout":
				next(w, r)
				return
			}

			cookie, err := r.Cookie("csrfToken")
			headerToken := r.Header.Get("X-CSRF-Token")
			if err != nil || cookie.Value == "" || headerToken == "" || cookie.Value != headerToken {
				http.Error(w, "CSRF validation failed", http.StatusForbidden)
				return
			}
		}

		next(w, r)
	}
}
