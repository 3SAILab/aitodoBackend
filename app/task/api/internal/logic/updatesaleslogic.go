package logic

import (
	"context"
	"errors"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesLogic {
	return &UpdateSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSalesLogic) UpdateSales(req *types.UpdateSalesReq) error {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}

	_, err := l.svcCtx.SalesPersonModel.FindOne(l.ctx, req.Id)
	if err != nil {
		return err
	}
	sales := &model.SalesPerson{
		Id:    req.Id,
		Name:  req.Name,
		Phone: req.Phone,
	}
	return l.svcCtx.SalesPersonModel.Update(l.ctx, sales)
}
