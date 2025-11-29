package handler

import (
	"net/http"
	"todo/app/task/api/internal/logic"
	"todo/app/task/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewListTaskLogic(r.Context(), svcCtx)
		resp, err := l.ListTask()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
