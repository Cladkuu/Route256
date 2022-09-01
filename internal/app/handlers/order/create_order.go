package order

import (
	"context"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) CreateOrder(ctx context.Context, request *desc.CreateOrderRequest) (*desc.CreateOrderResponse, error) {
	order, err := i.orderService.CreateOrder(ctx, request.GetPrice(),
		request.GetCurrency(),
		request.GetOrderCode())
	if err != nil {
		return nil, err
	}

	return &desc.CreateOrderResponse{Id: order}, nil
}
