package logic

import (
	"context"
	"database/sql"
	"time"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(req *types.CreateTaskReq) (resp *types.CreateTaskResp, err error) {
	// 1. 获取当前用户ID (假设 JWT 解析后的 key 是 "userId")
	userId := ""
	if v := l.ctx.Value("userId"); v != nil {
		userId = v.(string)
	}

	taskId := uuid.New().String()
	cstZone := time.FixedZone("CST", 8*3600) // 东八区
	now := time.Now().In(cstZone)

	assignee := sql.NullString{}
	if req.AssigneeId != "" {
		assignee = sql.NullString{String: req.AssigneeId, Valid: true}
	}

	sales := sql.NullString{}
	if req.SalesPersonId != "" {
		sales = sql.NullString{String: req.SalesPersonId, Valid: true}
	}

	desc := sql.NullString{}
	if req.Description != "" {
		desc = sql.NullString{String: req.Description, Valid: true}
	}

	// 处理日期
	due := sql.NullTime{}
	if req.DueDate > 0 {
		due = sql.NullTime{Time: time.Unix(req.DueDate, 0).UTC(), Valid: true}
	}

	newTask := &model.Task{
		Id:            taskId,
		TypeId:        req.TypeId,
		CreatorId:     userId,
		AssigneeId:    assignee,
		SalesPersonId: sales,
		Title:         req.Title,
		Description:   desc,
		Status:        "TODO", // 默认状态
		SortOrder:     int64(req.Priority),
		DueDate:       due,
		IsDeleted:     false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	_, err = l.svcCtx.TaskModel.Insert(l.ctx, newTask)
	if err != nil {
		return nil, err
	}
	return &types.CreateTaskResp{Id: taskId}, nil
}
