package logic

import (
	"context"
	"errors"
	"time"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"
	"todo/app/task/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSalesLogic {
	return &CreateSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSalesLogic) CreateSales(req *types.CreateSalesReq) (resp *types.SalesResp, err error) {
	// 统一添加的权限校验代码
	role := ""
	if v := l.ctx.Value("role"); v != nil {
		role = v.(string)
	}
	if role != "admin" {
		return nil, errors.New("权限不足：非管理员禁止操作") // 或者 return errors.New(...)
	}
	id := uuid.New().String()
	now := time.Now()

	sales := &model.SalesPerson{
		Id:        id,
		Name:      req.Name,
		Phone:     req.Phone,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = l.svcCtx.SalesPersonModel.Insert(l.ctx, sales)
	if err != nil {
		return nil, err
	}
	return &types.SalesResp{
		Id:    sales.Id,
		Name:  sales.Name,
		Phone: sales.Phone,
	}, nil
}
