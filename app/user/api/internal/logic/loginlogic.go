package logic

import (
	"context"
	"errors"
	"time"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"
	"todo/app/user/model"

	"github.com/golang-jwt/jwt/v4"
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

	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, accessExpire, user.Id, user.Role)
	if err != nil {
		return nil, err
	}
	return &types.LoginResp{
		AccessToken:  accessToken,
		AccessExpire: accessExpire,
		Id:           user.Id,
		Username:     user.Username,
		Role:         user.Role,
	}, nil
}

func (l *LoginLogic) getJwtToken(secretKey string, iat, seconds int64, userId string, userRole string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["role"] = userRole // [新增] 将角色放入 Token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
