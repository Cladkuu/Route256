package order

import (
	"context"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/convertors"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) GetAllOrders(ctx context.Context, req *desc.GetAllOrdersRequest) (*desc.GetAllOrdersResponse, error) {

	orders, err := i.orderService.GetAllOrders(ctx,
		req.GetPagination().GetPage(),
		req.GetPagination().GetPageSize(),
		req.GetPagination().GetSort().String())
	if err != nil {
		return nil, err
	}

	return convertors.ConvertToAllOrders(orders), nil

}
