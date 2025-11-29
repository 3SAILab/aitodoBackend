package logic

import (
	"context"
	"errors"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTaskTypeLogic {
	return &DeleteTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTaskTypeLogic) DeleteTaskType(req *types.DeleteTaskTypeReq) error {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}
	// 这里可以加一个逻辑：如果该类型下还有任务，禁止删除
	// 但为了简单起见，这里直接执行删除
	count, err := l.svcCtx.TaskModel.CountByTypeId(l.ctx, req.Id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该类型下仍有正在进行的任务，禁止删除")
	}
	return l.svcCtx.TaskTypeModel.Delete(l.ctx, req.Id)
}
