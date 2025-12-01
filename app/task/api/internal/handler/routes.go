package handler

import (
	"net/http"
	"todo/app/task/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/tasks",
				Handler: CreateTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/tasks/:id",
				Handler: DeleteTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/sales",
				Handler: CreateSalesHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/sales",
				Handler: ListSalesHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/sales/:id",
				Handler: UpdateSalesHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/sales/:id",
				Handler: DeleteSalesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/task-types",
				Handler: CreateTaskTypeHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/task-types",
				Handler: ListTaskTypeHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/task-types/:id",
				Handler: UpdateTaskTypeHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/task-types/:id",
				Handler: DeleteTaskTypeHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tasks",
				Handler: ListTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/tasks/:id",
				Handler: UpdateTaskHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/tasks/:taskId/progress",
				Handler: CreateTaskProgressHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/tasks/:taskId/progress",
				Handler: ListTaskProgressHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
