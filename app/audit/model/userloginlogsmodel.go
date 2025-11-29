package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserLoginLogsModel = (*customUserLoginLogsModel)(nil)

type (
	// UserLoginLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginLogsModel.
	UserLoginLogsModel interface {
		userLoginLogsModel
		withSession(session sqlx.Session) UserLoginLogsModel
	}

	customUserLoginLogsModel struct {
		*defaultUserLoginLogsModel
	}
)

// NewUserLoginLogsModel returns a model for the database table.
func NewUserLoginLogsModel(conn sqlx.SqlConn) UserLoginLogsModel {
	return &customUserLoginLogsModel{
		defaultUserLoginLogsModel: newUserLoginLogsModel(conn),
	}
}

func (m *customUserLoginLogsModel) withSession(session sqlx.Session) UserLoginLogsModel {
	return NewUserLoginLogsModel(sqlx.NewSqlConnFromSession(session))
}
