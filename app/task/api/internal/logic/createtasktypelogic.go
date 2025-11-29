package logic

import (
	"context"
	"errors"
	"time"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskTypeLogic {
	return &CreateTaskTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskTypeLogic) CreateTaskType(req *types.CreateTaskTypeReq) (resp *types.TaskTypeResp, err error) {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return nil, errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}
	// 获取当前操作用户 ID
	userId := ""
	if v := l.ctx.Value("userId"); v != nil {
		userId = v.(string)
	}
	if userId == "" {
		return nil, errors.New("无法获取当前用户身份")
	}

	id := uuid.New().String()
	now := time.Now()

	taskType := &model.TaskType{
		Id:        id,
		Name:      req.Name,
		ColorCode: req.ColorCode,
		CreatedBy: userId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = l.svcCtx.TaskTypeModel.Insert(l.ctx, taskType)
	if err != nil {
		return nil, err
	}

	return &types.TaskTypeResp{
		Id:        taskType.Id,
		Name:      taskType.Name,
		ColorCode: taskType.ColorCode,
	}, nil
}
