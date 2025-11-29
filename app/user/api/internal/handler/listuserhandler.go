package handler

import (
	"net/http"
	"todo/app/user/api/internal/logic"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListUserLogic(r.Context(), svcCtx)
		resp, err := l.ListUser()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
