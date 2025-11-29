package logic

import (
	"context"
	"errors"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"
	"todo/app/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) error {
	userIdFromCtx := l.ctx.Value("userId")
	if userIdFromCtx == nil {
		return errors.New("未授权：无法获取用户信息")
	}
	currentUserId, ok := userIdFromCtx.(string) // 这里对应你在 loginlogic 里存入的类型
	if !ok {
		return errors.New("用户身份解析失败")
	}

	currentUser, err := l.svcCtx.UserModel.FindOne(l.ctx, currentUserId)
	if err != nil {
		if err == model.ErrNotFound {
			return errors.New("操作者账户不存在或已被删除")
		}
		return err
	}

	if currentUser.Role != "admin" {
		return errors.New("权限不足：只有管理员可以执行删除操作")
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return err
	}
	if user.Role == "admin" {
		return errors.New("禁止删除管理员账户")
	}
	err = l.svcCtx.UserModel.Delete(l.ctx, req.Id)
	if err != nil {
		return err
	}
	return nil
}
