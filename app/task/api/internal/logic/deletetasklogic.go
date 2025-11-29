package logic

import (
	"context"
	"errors"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskLogic {
	return &DeleteTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskLogic) DeleteTask(req *types.DeleteTaskReq) error {
	userIdFromCtx := l.ctx.Value("userId")
	if userIdFromCtx == nil {
		return errors.New("未授权")
	}
	currentUserId, ok := userIdFromCtx.(string)
	if !ok {
		return errors.New("身份解析错误")
	}
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return errors.New("任务不存在")
		}
		return err
	}

	if task.IsDeleted {
		return errors.New("任务不存在")
	}

	if task.CreatorId != currentUserId {
		return errors.New("无权删除该任务：只有创建人可以删除")
	}
	err = l.svcCtx.TaskModel.Delete(l.ctx, req.Id)
	return err
}
