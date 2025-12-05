package logic

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"
	"todo/app/user/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成新的 sessionId，并写入 Redis，实现单设备登录与 refreshToken 版本管理
	// 开发环境下如果 Redis 未启动，不阻断登录，只记录日志
	sessionId := uuid.New().String()
	sessionKey := fmt.Sprintf("user:session:%s", user.Id)
	ttl := int(l.svcCtx.Config.Auth.AccessExpire * 7)
	if err := l.svcCtx.Redis.SetexCtx(l.ctx, sessionKey, sessionId, ttl); err != nil {
		l.Logger.Errorf("failed to write session to redis, continue without redis session: %v", err)
	}

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire

	// accessToken 携带 sid，方便后端识别当前会话
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.Id, user.Role, sessionId)
	if err != nil {
		return nil, err
	}

	// refreshToken 有更长有效期，使用相同 sid。这里简单用 7 倍 access 过期时间
	refreshExpire := accessExpire * 7
	refreshToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, refreshExpire, user.Id, user.Role, sessionId)
	if err != nil {
		return nil, err
	}

	// 生成 CSRF Token（双重提交：Cookie + Header）
	csrfToken := uuid.New().String()

	return &types.LoginResp{
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
		Id:           user.Id,
		Username:     user.Username,
		Role:         user.Role,
		SessionId:    sessionId,
		CsrfToken:    csrfToken,
		RefreshToken: refreshToken,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds int64, userId string, userRole string, sessionId string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["role"] = userRole // [新增] 将角色放入 Token
	claims["sid"] = sessionId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
