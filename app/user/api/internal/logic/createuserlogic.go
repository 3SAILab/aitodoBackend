package logic

import (
	"context"
	"errors"
	"time"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"
	"todo/app/user/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) (resp *types.UserResp, err error) {
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string) // 这里对应你在 LoginLogic 里存入的 key
	}
	if role != "admin" {
		return nil, errors.New("权限不足：只有管理员可以创建新用户")
	}

	_, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, req.Email)
	if err == nil {
		return nil, errors.New("用户已存在")
	}
	if err != model.ErrNotFound {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newId := uuid.New().String()
	user := &model.User{
		Id:           newId,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &types.UserResp{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, err
}
