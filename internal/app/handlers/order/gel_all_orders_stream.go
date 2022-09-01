package order

import (
	"database/sql"
	"github.com/pkg/errors"
	desc "gitlab.ozon.dev/astoyakin/route256/pkg/order-api"
	"io"
)

func (i *Implementation) GetAllOrdersStream(stream desc.OrderApi_GetAllOrdersStreamServer) error {
	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		order, err := i.orderService.GetOrderById(ctx, req.GetId())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		if err = stream.Send(&desc.GetOrderByIdResponse{Order: &desc.Order{
			Currency: order.GetCurrency(),
			Status:   order.GetStatus(),
			Price:    order.GetPrice(),
		}}); err != nil {
			return err
		}
	}

}
