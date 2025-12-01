package logic

import (
	"context"
	"log"
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
	log.Printf("TaskId=%s | ListByTaskId 结果：list长度=%d, list内容=%+v", req.TaskId, len(list), list)
	for i, item := range list {
		// %+v 会显示结构体的「字段名+值」，即使是指针也能解析
		log.Printf("TaskId=%s | list[%d] 内容：%+v", req.TaskId, i, item)
	}
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
