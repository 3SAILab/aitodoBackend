package logic

import (
	"context"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskTypeLogic {
	return &ListTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskTypeLogic) ListTaskType() (resp *types.ListTaskTypeResp, err error) {
	list, err := l.svcCtx.TaskTypeModel.List(l.ctx)
	if err != nil {
		return nil, err
	}

	var respList []types.TaskTypeResp
	for _, item := range list {
		respList = append(respList, types.TaskTypeResp{
			Id:        item.Id,
			Name:      item.Name,
			ColorCode: item.ColorCode,
		})
	}

	return &types.ListTaskTypeResp{
		List: respList,
	}, nil
}
