package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserBehaviorLogsModel = (*customUserBehaviorLogsModel)(nil)

type (
	// UserBehaviorLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserBehaviorLogsModel.
	UserBehaviorLogsModel interface {
		userBehaviorLogsModel
		withSession(session sqlx.Session) UserBehaviorLogsModel
	}

	customUserBehaviorLogsModel struct {
		*defaultUserBehaviorLogsModel
	}
)

// NewUserBehaviorLogsModel returns a model for the database table.
func NewUserBehaviorLogsModel(conn sqlx.SqlConn) UserBehaviorLogsModel {
	return &customUserBehaviorLogsModel{
		defaultUserBehaviorLogsModel: newUserBehaviorLogsModel(conn),
	}
}

func (m *customUserBehaviorLogsModel) withSession(session sqlx.Session) UserBehaviorLogsModel {
	return NewUserBehaviorLogsModel(sqlx.NewSqlConnFromSession(session))
}
