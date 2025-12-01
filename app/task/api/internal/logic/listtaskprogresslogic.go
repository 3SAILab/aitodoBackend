package logic

import (
	"context"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskProgressLogic {
	return &ListTaskProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskProgressLogic) ListTaskProgress(req *types.ListTaskProgressReq) (resp *types.ListTaskProgressResp, err error) {
	list, err := l.svcCtx.TaskProgressModel.ListByTaskId(l.ctx, req.TaskId)
	if err != nil {
		return nil, err
	}
	var respList []types.TaskProgressResp
	for _, item := range list {
		respList = append(respList, types.TaskProgressResp{
			Id:        item.Id,
			TaskId:    item.TaskId,
			Content:   item.Content,
			CreatedBy: item.CreatedBy,
			CreatedAt: item.CreatedAt.Unix(),
		})
	}
	return &types.ListTaskProgressResp{List: respList}, nil
}
