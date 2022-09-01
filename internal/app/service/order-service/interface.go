package order_service

import (
	"context"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
)

type IOrder interface {
	CreateOrder(ctx context.Context, price int32, currency string, orderCode string) (int64, error)
	GetAllOrders(ctx context.Context, page, pageSize int32, sortCriteria string) ([]*model.Order, error)
	GetOrderById(ctx context.Context, id int64) (*model.Order, error)
	CancelOrder(ctx context.Context, id int64) error
	ChangeOrderStatus(ctx context.Context, id int64, statusTo string) error
	ResetOrderPrice(ctx context.Context, id int64) error
}
