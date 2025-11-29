package logic

import (
	"context"
	"errors"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSalesLogic {
	return &DeleteSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSalesLogic) DeleteSales(req *types.DeleteSalesReq) error {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}
	return l.svcCtx.SalesPersonModel.Delete(l.ctx, req.Id)
}
