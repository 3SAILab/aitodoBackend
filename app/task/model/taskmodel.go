package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Task struct {
	Id            string         `db:"id"`
	TypeId        string         `db:"type_id"`
	CreatorId     string         `db:"creator_id"`
	AssigneeId    sql.NullString `db:"assignee_id"`
	SalesPersonId sql.NullString `db:"sales_person_id"`
	Title         string         `db:"title"`
	Description   sql.NullString `db:"description"`
	Status        string         `db:"status"`
	SortOrder     int64          `db:"sort_order"`
	DueDate       sql.NullTime   `db:"due_date"`
	IsDeleted     bool           `db:"is_deleted"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
	CompletedAt   sql.NullTime   `db:"completed_at"`
}

type TasksModel interface {
	Insert(ctx context.Context, task *Task) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*Task, error)
	Delete(ctx context.Context, id string) error
	CountByTypeId(ctx context.Context, typeId string) (int64, error)
	Update(ctx context.Context, task *Task) error
	List(ctx context.Context) ([]*Task, error)
}

type defaultTasksModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewTasksModel(conn sqlx.SqlConn) TasksModel {
	return &defaultTasksModel{
		conn:  conn,
		table: `"public"."tasks"`, // PostgreSQL 表名通常加引号
	}
}

func (m *defaultTasksModel) Insert(ctx context.Context, task *Task) (sql.Result, error) {
	query := `INSERT INTO "public"."tasks" 
		(id, type_id, creator_id, assignee_id, sales_person_id, title, description, status, sort_order, due_date, is_deleted, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	return m.conn.ExecCtx(ctx, query,
		task.Id, task.TypeId, task.CreatorId, task.AssigneeId, task.SalesPersonId,
		task.Title, task.Description, task.Status, task.SortOrder, task.DueDate,
		task.IsDeleted, task.CreatedAt, task.UpdatedAt,
	)
}

func (m *defaultTasksModel) FindOne(ctx context.Context, id string) (*Task, error) {
	query := `SELECT id, type_id, creator_id, assignee_id, sales_person_id, title, description, status, sort_order, due_date, is_deleted, created_at, updated_at, completed_at 
		FROM "public"."tasks" WHERE id = $1 LIMIT 1`
	var resp Task
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

func (m *defaultTasksModel) Delete(ctx context.Context, id string) error {
	query := `UPDATE "public"."tasks" SET is_deleted = true, updated_at = $1 WHERE id = $2`
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}

func (m *defaultTasksModel) CountByTypeId(ctx context.Context, typeId string) (int64, error) {
	// 只统计 is_deleted = false 的任务
	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE type_id = $1 AND is_deleted = false`, m.table)
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, query, typeId)
	return count, err
}

func (m *defaultTasksModel) Update(ctx context.Context, task *Task) error {
	query := `UPDATE "public"."tasks" 
		SET type_id = $1, assignee_id = $2, sales_person_id = $3, title = $4, 
		description = $5, status = $6, sort_order = $7, due_date = $8, 
		completed_at = $9, updated_at = $10 
		WHERE id = $11`

	_, err := m.conn.ExecCtx(ctx, query,
		task.TypeId, task.AssigneeId, task.SalesPersonId, task.Title,
		task.Description, task.Status, task.SortOrder, task.DueDate,
		task.CompletedAt, // [修改点 2] 传入完成时间
		time.Now(),       // updated_at
		task.Id,          // WHERE id
	)
	return err
}

func (m *defaultTasksModel) List(ctx context.Context) ([]*Task, error) {
	query := `SELECT id, type_id, creator_id, assignee_id, sales_person_id, title, description, status, sort_order, due_date, is_deleted, created_at, updated_at, completed_at 
		FROM "public"."tasks" 
		WHERE is_deleted = false 
		ORDER BY sort_order DESC, created_at DESC`
	var resp []*Task
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
