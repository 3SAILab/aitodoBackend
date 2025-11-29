package logic

import (
	"context"
	"todo/app/task/api/internal/svc"
	"todo/app/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSalesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSalesLogic {
	return &ListSalesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSalesLogic) ListSales() (resp *types.ListSalesResp, err error) {
	list, err := l.svcCtx.SalesPersonModel.List(l.ctx)
	if err != nil {
		return nil, err
	}
	var respList []types.SalesResp
	for _, item := range list {
		respList = append(respList, types.SalesResp{
			Id:    item.Id,
			Name:  item.Name,
			Phone: item.Phone,
		})
	}
	return &types.ListSalesResp{
		List: respList,
	}, nil
}
