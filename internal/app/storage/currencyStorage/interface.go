package currencyStorage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type ICurrency interface {
	GetCurrency(ctx context.Context, transaction pgx.Tx, currency string) error
}
