package order

import (
	"context"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.CancelOrderRequest) (*desc.CancelOrderResponse, error) {
	if err := i.orderService.CancelOrder(ctx,
		req.GetId()); err != nil {
		return nil, err
	}

	return &desc.CancelOrderResponse{Result: true}, nil
}
