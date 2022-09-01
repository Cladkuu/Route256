package orderStorage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
)

type IOrderStorage interface {
	GetAllOrders(ctx context.Context, transaction pgx.Tx, page, pageSize int32, sortCriteria string) ([]*model.Order, error)
	GetOrderById(ctx context.Context, transaction pgx.Tx, id int64) (*model.Order, error)
	CreateOrder(ctx context.Context, transaction pgx.Tx, price int32, currency, orderCode string) (int64, error)
	ChangeStatus(ctx context.Context, transaction pgx.Tx, id int64, status string) error
	CancelOrder(ctx context.Context, transaction pgx.Tx, id int64) error
	ResetOrderPrice(ctx context.Context, transaction pgx.Tx, id int64) error
	GetOrderByOrderCode(ctx context.Context, transaction pgx.Tx, orderCode string) (*model.Order, error)
}
