package handler

import (
	"net/http"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/users/login",
				Handler: LoginHandler(serverCtx),
			},
		},
	)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/users",
				Handler: CreateUserHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/users/:id",
				Handler: DeleteUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/users",
				Handler: ListUserHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
