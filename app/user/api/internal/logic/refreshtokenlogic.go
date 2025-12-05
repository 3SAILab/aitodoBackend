package logic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"
	"todo/app/user/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// RefreshTokenLogic 负责根据 HttpOnly refreshToken Cookie 刷新 accessToken。
// 说明：此处为了示例简单，使用与 accessToken 相同的签名密钥与结构；
// 生产环境建议为 refreshToken 使用独立密钥、更长有效期，并在服务端维护 refreshToken 黑名单或版本号进行 rotation 与撤销控制。
type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(w http.ResponseWriter, r *http.Request) (resp *types.RefreshTokenResp, err error) {
	refreshCookie, err := r.Cookie("refreshToken")
	if err != nil || refreshCookie.Value == "" {
		return nil, errors.New("missing refresh token")
	}

	refreshToken := refreshCookie.Value

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	userId, _ := claims["userId"].(string)
	role, _ := claims["role"].(string)
	oldSid, _ := claims["sid"].(string)

	// 从 Redis 读取当前有效 sessionId
	sessionKey := fmt.Sprintf("user:session:%s", userId)
	currentSid, err := l.svcCtx.Redis.GetCtx(l.ctx, sessionKey)
	if err != nil {
		// 开发环境可能未启动 Redis：此时放宽校验，仅记录日志
		l.Logger.Errorf("failed to read session from redis, skip validation: %v", err)
		currentSid = oldSid
	}
	if currentSid == "" {
		currentSid = oldSid
	}
	if currentSid != oldSid {
		return nil, errors.New("refresh token revoked")
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New("user is disabled")
	}

	// 生成新的 sid，实现 rotation，并覆盖 Redis
	newSid := uuid.New().String()
	ttl := int(l.svcCtx.Config.Auth.AccessExpire * 7)
	if err := l.svcCtx.Redis.SetexCtx(l.ctx, sessionKey, newSid, ttl); err != nil {
		l.Logger.Errorf("failed to write session to redis during refresh, continue without redis session: %v", err)
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.Id, role, newSid)
	if err != nil {
		return nil, err
	}

	// 刷新 refreshToken，使用新的 sid
	refreshExpire := accessExpire * 7
	newRefreshToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, refreshExpire, user.Id, role, newSid)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 生产环境必须使用 HTTPS，并设置为 true
		SameSite: http.SameSiteLaxMode,
	})

	return &types.RefreshTokenResp{
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
		Id:           user.Id,
		Username:     user.Username,
		Role:         user.Role,
		SessionId:    newSid,
	}, nil
}

// 复用 LoginLogic 中的生成 Token 逻辑
func (l *RefreshTokenLogic) getJwtToken(secretKey string, iat, seconds int64, userId string, userRole string, sessionId string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["role"] = userRole
	claims["sid"] = sessionId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
