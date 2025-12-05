package handler

import (
	"net/http"
	"todo/app/user/api/internal/logic"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 登录成功：将 refreshToken 写入 HttpOnly Cookie（前端无法读取，仅用于刷新）
		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    resp.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // 生产环境必须开启 HTTPS 并置为 true
			SameSite: http.SameSiteLaxMode,
		})

		// 下发 CSRF Token：使用非 HttpOnly Cookie + JSON 字段，配合前端 Header 做双重提交校验
		if resp.CsrfToken != "" {
			http.SetCookie(w, &http.Cookie{
				Name:     "csrfToken",
				Value:    resp.CsrfToken,
				Path:     "/",
				HttpOnly: false, // 需要前端 JS 可读取时必须为 false
				Secure:   false, // 生产环境下应为 true 并启用 HTTPS
				SameSite: http.SameSiteLaxMode,
			})
		}

		// 返回给前端的 JSON 中不会包含 RefreshToken 字段（struct tag 已经排除）
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
