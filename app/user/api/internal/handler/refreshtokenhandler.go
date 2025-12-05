package handler

import (
	"net/http"
	"todo/app/user/api/internal/logic"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// RefreshTokenHandler 根据 HttpOnly refreshToken Cookie 刷新 accessToken。
func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}


