package logic

import (
	"context"
	"errors"
	"time"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"
	userModel "todo/app/user/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskProgressLogic {
	return &CreateTaskProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskProgressLogic) CreateTaskProgress(req *types.CreateTaskProgressReq) (resp *types.TaskProgressResp, err error) {
	userId := ""
	if v := l.ctx.Value("userId"); v != nil {
		userId = v.(string)
	}
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.TaskId)
	if err != nil {
		if err == userModel.ErrNotFound {
			return nil, errors.New("任务不存在，无法添加进度")
		}
		return nil, err
	}
	if task.IsDeleted {
		return nil, errors.New("任务已删除，无法添加进度")
	}
	id := uuid.New().String()
	cstZone := time.FixedZone("CST", 8*3600) // 东八区
	now := time.Now().In(cstZone)

	data := &model.TaskProgressLog{
		Id:        id,
		TaskId:    req.TaskId,
		Content:   req.Content,
		CreatedBy: userId,
		CreatedAt: now,
	}
	_, err = l.svcCtx.TaskProgressModel.Insert(l.ctx, data)
	if err != nil {
		return nil, err
	}
	return &types.TaskProgressResp{
		Id:        id,
		TaskId:    req.TaskId,
		Content:   req.Content,
		CreatedBy: userId,
		CreatedAt: now.Unix(),
	}, nil
}
