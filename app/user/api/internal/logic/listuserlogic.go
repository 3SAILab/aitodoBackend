package logic

import (
	"context"
	"todo/app/user/api/internal/svc"
	"todo/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLogic) ListUser() (*types.ListUserResp, error) {

	users, err := l.svcCtx.UserModel.List(l.ctx)
	if err != nil {
		return nil, err
	}

	var respList []types.UserResp
	for _, t := range users {
		respList = append(respList, types.UserResp{
			Id:       t.Id,
			Username: t.Username, // 用户显示名称
			Email:    t.Email,    // 登录邮箱（唯一）
			Role:     t.Role,     // 角色（admin/user）
		})
	}

	return &types.ListUserResp{
		List: respList,
	}, nil
}
