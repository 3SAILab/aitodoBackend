package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TaskActivityLogsModel = (*customTaskActivityLogsModel)(nil)

type (
	// TaskActivityLogsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskActivityLogsModel.
	TaskActivityLogsModel interface {
		taskActivityLogsModel
		withSession(session sqlx.Session) TaskActivityLogsModel
	}

	customTaskActivityLogsModel struct {
		*defaultTaskActivityLogsModel
	}
)

// NewTaskActivityLogsModel returns a model for the database table.
func NewTaskActivityLogsModel(conn sqlx.SqlConn) TaskActivityLogsModel {
	return &customTaskActivityLogsModel{
		defaultTaskActivityLogsModel: newTaskActivityLogsModel(conn),
	}
}

func (m *customTaskActivityLogsModel) withSession(session sqlx.Session) TaskActivityLogsModel {
	return NewTaskActivityLogsModel(sqlx.NewSqlConnFromSession(session))
}
