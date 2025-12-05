package logic

import (
	"context"
	"fmt"
	"net/http"
	"todo/app/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// LogoutLogic 用于清理客户端 refreshToken Cookie。
// 说明：实际的 refreshToken 撤销应在服务端维护黑名单或版本号，这里示例只删除 Cookie。
type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(userId string, w http.ResponseWriter) error {
	// 删除 Redis 中的会话信息，实现强制下线
	if userId != "" {
		sessionKey := fmt.Sprintf("user:session:%s", userId)
		_, _ = l.svcCtx.Redis.DelCtx(l.ctx, sessionKey)
	}

	// 通过设置过期时间清除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // 生产环境必须使用 HTTPS，并设置为 true
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}


