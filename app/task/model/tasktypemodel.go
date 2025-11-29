package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type TaskType struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	ColorCode string    `db:"color_code"`
	CreatedBy string    `db:"created_by"` // 关联 users.id
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type TaskTypeModel interface {
	Insert(ctx context.Context, data *TaskType) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*TaskType, error)
	List(ctx context.Context) ([]*TaskType, error)
	Update(ctx context.Context, data *TaskType) error
	Delete(ctx context.Context, id string) error
}

type defaultTaskTypeModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewTaskTypeModel(conn sqlx.SqlConn) TaskTypeModel {
	return &defaultTaskTypeModel{
		conn:  conn,
		table: `"public"."task_types"`,
	}
}

func (m *defaultTaskTypeModel) Insert(ctx context.Context, data *TaskType) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, name, color_code, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, m.table)
	return m.conn.ExecCtx(ctx, query, data.Id, data.Name, data.ColorCode, data.CreatedBy, data.CreatedAt, data.UpdatedAt)
}

func (m *defaultTaskTypeModel) FindOne(ctx context.Context, id string) (*TaskType, error) {
	query := fmt.Sprintf(`SELECT id, name, color_code, created_by, created_at, updated_at FROM %s WHERE id = $1 LIMIT 1`, m.table)
	var resp TaskType
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, sqlx.ErrNotFound
	default:
		return nil, err
	}
}

// List 列出所有任务类型
func (m *defaultTaskTypeModel) List(ctx context.Context) ([]*TaskType, error) {
	query := fmt.Sprintf(`SELECT id, name, color_code, created_by, created_at, updated_at FROM %s ORDER BY created_at ASC`, m.table)
	var resp []*TaskType
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

func (m *defaultTaskTypeModel) Update(ctx context.Context, data *TaskType) error {
	query := fmt.Sprintf(`UPDATE %s SET name = $1, color_code = $2, updated_at = $3 WHERE id = $4`, m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.ColorCode, time.Now(), data.Id)
	return err
}

// Delete 物理删除（因为表结构中没有 is_deleted 字段）
func (m *defaultTaskTypeModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
