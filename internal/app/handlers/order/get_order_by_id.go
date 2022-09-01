package order

import (
	"context"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) GetOrderById(ctx context.Context, req *desc.GetOrderByIdRequest) (*desc.GetOrderByIdResponse, error) {
	order, err := i.orderService.GetOrderById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetOrderByIdResponse{Order: &desc.Order{
		Currency: order.GetCurrency(),
		Status:   order.GetStatus(),
		Price:    order.GetPrice(),
	}}, nil
}
