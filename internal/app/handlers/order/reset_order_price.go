package order

import (
	"context"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) ResetOrderPrice(ctx context.Context, req *desc.ResetOrderPriceRequest) (*desc.ResetOrderPriceResponse, error) {
	if err := i.orderService.ResetOrderPrice(ctx,
		req.GetId()); err != nil {
		return nil, err
	}

	return &desc.ResetOrderPriceResponse{Result: true}, nil
}
