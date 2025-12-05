package handler

import (
	"net/http"
	"strings"
	"todo/app/user/api/internal/logic"
	"todo/app/user/api/internal/svc"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// LogoutHandler 清理 refreshToken Cookie。
func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Authorization 头中解析出当前用户 ID，用于删除 Redis 会话
		userId := ""
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString := parts[1]
				claims := jwt.MapClaims{}
				token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(svcCtx.Config.Auth.AccessSecret), nil
				})
				if token != nil && token.Valid {
					if v, ok := claims["userId"].(string); ok {
						userId = v
					}
				}
			}
		}

		l := logic.NewLogoutLogic(r.Context(), svcCtx)
		if err := l.Logout(userId, w); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.Ok(w)
	}
}


