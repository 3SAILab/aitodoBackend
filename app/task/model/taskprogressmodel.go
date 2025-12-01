package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type TaskProgressLog struct {
	Id        string    `db:"id"`
	TaskId    string    `db:"task_id"`
	Content   string    `db:"content"`
	CreatedBy string    `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

type TaskProgressModel interface {
	Insert(ctx context.Context, data *TaskProgressLog) (sql.Result, error)
	ListByTaskId(ctx context.Context, taskId string) ([]*TaskProgressLog, error)
}

type defaultTaskProgressModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewTaskProgressModel(conn sqlx.SqlConn) TaskProgressModel {
	return &defaultTaskProgressModel{
		conn:  conn,
		table: `"public"."task_progress_logs"`,
	}
}

func (m *defaultTaskProgressModel) Insert(ctx context.Context, data *TaskProgressLog) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, task_id, content, created_by, created_at) VALUES ($1, $2, $3, $4, $5)`, m.table)
	return m.conn.ExecCtx(ctx, query, data.Id, data.TaskId, data.Content, data.CreatedBy, data.CreatedAt)
}

func (m *defaultTaskProgressModel) ListByTaskId(ctx context.Context, taskId string) ([]*TaskProgressLog, error) {
	query := fmt.Sprintf(`SELECT id, task_id, content, created_by, created_at FROM %s WHERE task_id = $1 ORDER BY created_at DESC`, m.table)
	var resp []*TaskProgressLog
	err := m.conn.QueryRowsCtx(ctx, &resp, query, taskId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}
