package handler

import (
	"net/http"
	"strings"
	"todo/app/user/api/internal/logic"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// VerifyTokenHandler 对前端传来的 Authorization Bearer token 做一次显式校验。
func VerifyTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.OkJsonCtx(r.Context(), w, map[string]any{
				"valid": false,
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httpx.OkJsonCtx(r.Context(), w, map[string]any{
				"valid": false,
			})
			return
		}

		tokenString := parts[1]

		l := logic.NewVerifyTokenLogic(r.Context(), svcCtx)
		resp, err := l.VerifyToken(tokenString)
		if err != nil || !resp.Valid {
			httpx.OkJsonCtx(r.Context(), w, map[string]any{
				"valid": false,
			})
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}



