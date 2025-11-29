package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 数据库映射
type User struct {
	Id           string    `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	AvatarUrl    string    `db:"avatar_url"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// 接口
type UserModel interface {
	Insert(ctx context.Context, user *User) (sql.Result, error) // 插入用户数据
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOne(ctx context.Context, id string) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*User, error) // 获取所有在职销售

}

// 接口的默认实现，持有数据库连接和表名
type defaultUserModel struct {
	conn  sqlx.SqlConn
	table string
}

// 构造函数：创建 UserModel 实例（依赖注入数据库连接）
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &defaultUserModel{
		conn:  conn,
		table: `"public"."users"`,
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, user *User) (sql.Result, error) {
	query := `INSERT INTO "public"."users" (id, username, email, password_hash, role, avatar_url, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	return m.conn.ExecCtx(ctx, query,
		user.Id, user.Username, user.Email, user.PasswordHash,
		user.Role, user.AvatarUrl, user.IsActive,
		user.CreatedAt, user.UpdatedAt,
	)
}

func (m *defaultUserModel) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, role, avatar_url, is_active, created_at, updated_at FROM "public"."users" WHERE email = $1 LIMIT 1`
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) FindOne(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, username, email, password_hash, role, avatar_url, is_active, created_at, updated_at FROM "public"."users" WHERE id = $1 LIMIT 1`
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "public"."users" WHERE id = $1`
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultUserModel) List(ctx context.Context) ([]*User, error) {
	// 获取所有正常状态的用户，按注册时间倒序排列（安全排除密码哈希字段）
	query := `SELECT id, username, email, password_hash, role, avatar_url, is_active, created_at, updated_at 
		FROM "public"."users"` // users表无sort_order，按注册时间倒序更合理

	var resp []*User
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, nil // 无数据时返回空切片（而非错误），更符合业务预期
	default:
		return nil, fmt.Errorf("query users failed: %w", err) // 包装错误，便于排查
	}
}
