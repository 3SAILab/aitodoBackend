package logic

import (
	"context"
	"errors"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

// VerifyTokenLogic 用于在需要时对 accessToken 做一次显式校验
// 一般用于路由守卫或后端内部服务间调用。
type VerifyTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTokenLogic {
	return &VerifyTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyTokenLogic) VerifyToken(tokenString string) (*types.VerifyTokenResp, error) {
	if tokenString == "" {
		return &types.VerifyTokenResp{Valid: false}, errors.New("empty token")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svcCtx.Config.Auth.AccessSecret), nil
	})
	if err != nil || !token.Valid {
		return &types.VerifyTokenResp{Valid: false}, err
	}

	userId, _ := claims["userId"].(string)
	role, _ := claims["role"].(string)

	if userId == "" {
		return &types.VerifyTokenResp{Valid: false}, errors.New("invalid token payload")
	}

	return &types.VerifyTokenResp{
		Valid: true,
		User: &types.UserResp{
			Id:   userId,
			Role: role,
		},
	}, nil
}



