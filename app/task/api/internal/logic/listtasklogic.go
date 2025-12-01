package logic

import (
	"context"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskLogic) ListTask() (resp *types.ListTaskResp, err error) {
	tasks, err := l.svcCtx.TaskModel.List(l.ctx)
	if err != nil {
		return nil, err
	}

	var respList []types.TaskResp
	for _, t := range tasks {
		dueDate := int64(0)
		if t.DueDate.Valid {
			dueDate = t.DueDate.Time.Unix()
		}
		completedAt := int64(0)
		if t.CompletedAt.Valid {
			completedAt = t.CompletedAt.Time.Unix()
		}

		// fmt.Printf("数据库 created_at: %v\n", t.CreatedAt)       // 应该显示 2025-11-29 10:22:14 +0800 CST
		// fmt.Printf("CreatedAt 时间戳: %v\n", t.CreatedAt.Unix()) // 对应 1759160534（示例值）

		// if t.DueDate.Valid {
		// 	fmt.Printf("数据库 due_date: %v\n", t.DueDate.Time)       // 应该显示 2025-11-29 18:26:00 +0800 CST
		// 	fmt.Printf("DueDate 时间戳: %v\n", t.DueDate.Time.Unix()) // 对应 1759189560（示例值）
		// }

		respList = append(respList, types.TaskResp{
			Id:            t.Id,
			TypeId:        t.TypeId,
			CreatorId:     t.CreatorId,
			AssigneeId:    t.AssigneeId.String,
			SalesPersonId: t.SalesPersonId.String,
			Title:         t.Title,
			Description:   t.Description.String,
			Status:        t.Status,
			Priority:      int(t.SortOrder),
			DueDate:       dueDate,
			CreatedAt:     t.CreatedAt.Unix(),
			CompletedAt:   completedAt,
		})
	}
	// 比如在构造响应前，打印关键时间

	return &types.ListTaskResp{
		List: respList,
	}, nil
}
