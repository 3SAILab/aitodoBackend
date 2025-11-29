package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTaskLogic {
	return &UpdateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTaskLogic) UpdateTask(req *types.UpdateTaskReq) error {
	// 1. 检查任务是否存在
	task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return err
	}

	// [新增逻辑开始] 处理 CompletedAt
	// 假设 "DONE" 是表示完成的状态字符串，请根据你的实际业务调整
	targetStatus := req.Status

	// 如果目标状态是完成，且当前状态不是完成 -> 设置当前时间
	if targetStatus == "DONE" && task.Status != "DONE" {
		task.CompletedAt = sql.NullTime{Time: time.Now(), Valid: true}
	} else if targetStatus != "DONE" && task.Status == "DONE" {
		// 如果从完成状态变更为未完成（例如重开任务）-> 清空完成时间
		task.CompletedAt = sql.NullTime{Valid: false}
	}
	// [新增逻辑结束]

	// 2. 准备更新数据 (处理 NullString 等类型)
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

	due := sql.NullTime{}
	if req.DueDate > 0 {
		due = sql.NullTime{Time: time.Unix(req.DueDate, 0), Valid: true}
		fmt.Println("截止时间: ", due)
	}

	// 3. 更新字段
	task.TypeId = req.TypeId
	task.Title = req.Title
	task.Description = desc
	task.AssigneeId = assignee
	task.SalesPersonId = sales
	task.Status = req.Status
	task.SortOrder = int64(req.Priority)
	task.DueDate = due
	// task.CompletedAt 已经在上面处理过了

	return l.svcCtx.TaskModel.Update(l.ctx, task)
}
