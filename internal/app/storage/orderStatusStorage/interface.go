package orderStatusStorage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type IOrderStatusStorage interface {
	FindStatuses(statusFrom, statusTo string) (bool, error)
	ValidateStatuses(ctx context.Context, transaction pgx.Tx, statusFrom, statusTo string) error
	ValidateStatusesV2(ctx context.Context, transaction pgx.Tx, id int64, statusTo string) error
}
