package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// SalesPerson 对应数据库表结构
type SalesPerson struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	Phone     string    `db:"phone"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SalesPersonModel interface {
	Insert(ctx context.Context, data *SalesPerson) (sql.Result, error)
	FindOne(ctx context.Context, id string) (*SalesPerson, error)
	List(ctx context.Context) ([]*SalesPerson, error) // 获取所有在职销售
	Update(ctx context.Context, data *SalesPerson) error
	Delete(ctx context.Context, id string) error // 软删除
}

type defaultSalesPersonModel struct {
	conn  sqlx.SqlConn
	table string
}

func NewSalesPersonModel(conn sqlx.SqlConn) SalesPersonModel {
	return &defaultSalesPersonModel{
		conn:  conn,
		table: `"public"."sales_persons"`, // PostgreSQL 表名
	}
}

func (m *defaultSalesPersonModel) Insert(ctx context.Context, data *SalesPerson) (sql.Result, error) {
	query := fmt.Sprintf(`INSERT INTO %s (id, name, phone, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`, m.table)
	return m.conn.ExecCtx(ctx, query, data.Id, data.Name, data.Phone, data.IsActive, data.CreatedAt, data.UpdatedAt)
}

func (m *defaultSalesPersonModel) FindOne(ctx context.Context, id string) (*SalesPerson, error) {
	query := fmt.Sprintf(`SELECT id, name, phone, is_active, created_at, updated_at FROM %s WHERE id = $1 LIMIT 1`, m.table)
	var resp SalesPerson
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

// List 列出所有“在职” (is_active=true) 的销售人员
func (m *defaultSalesPersonModel) List(ctx context.Context) ([]*SalesPerson, error) {
	query := fmt.Sprintf(`SELECT id, name, phone, is_active, created_at, updated_at FROM %s WHERE is_active = true ORDER BY created_at DESC`, m.table)
	var resp []*SalesPerson
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

func (m *defaultSalesPersonModel) Update(ctx context.Context, data *SalesPerson) error {
	query := fmt.Sprintf(`UPDATE %s SET name = $1, phone = $2, updated_at = $3 WHERE id = $4`, m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.Phone, time.Now(), data.Id)
	return err
}

// Delete 执行软删除，将 is_active 设为 false
func (m *defaultSalesPersonModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET is_active = false, updated_at = $1 WHERE id = $2`, m.table)
	_, err := m.conn.ExecCtx(ctx, query, time.Now(), id)
	return err
}
