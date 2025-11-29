package logic

import (
	"context"
	"errors"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskTypeLogic {
	return &UpdateTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskTypeLogic) UpdateTaskType(req *types.UpdateTaskTypeReq) error {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}

	// 检查是否存在
	_, err := l.svcCtx.TaskTypeModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return err
	}

	data := &model.TaskType{
		Id:        req.Id,
		Name:      req.Name,
		ColorCode: req.ColorCode,
	}
	return l.svcCtx.TaskTypeModel.Update(l.ctx, data)
}
