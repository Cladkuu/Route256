package order

import (
	"context"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
)

func (i *Implementation) ChangeStatus(ctx context.Context, req *desc.ChangeStatusRequest) (*desc.ChangeStatusResponse, error) {
	if err := i.orderService.ChangeOrderStatus(ctx,
		req.GetId(),
		req.Status.String()); err != nil {
		return nil, err
	}

	return &desc.ChangeStatusResponse{Result: true}, nil
}
