package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type IRepository interface {
	GetTransaction(ctx context.Context, TxOptions pgx.TxOptions) (pgx.Tx, context.Context, error)
	CommitTransaction(ctx context.Context, transaction pgx.Tx) error
	RollbackTransaction(ctx context.Context, transaction pgx.Tx) error
}
