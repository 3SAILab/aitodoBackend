package svc

import (
	"todo/app/task/api/internal/config"
	"todo/app/task/model"

	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	TaskModel        model.TasksModel
	SalesPersonModel model.SalesPersonModel
	TaskTypeModel    model.TaskTypeModel // [新增]
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)
	return &ServiceContext{
		Config:           c,
		TaskModel:        model.NewTasksModel(conn),
		SalesPersonModel: model.NewSalesPersonModel(conn),
		TaskTypeModel:    model.NewTaskTypeModel(conn), // [新增] 初始化
	}
}
