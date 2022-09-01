package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"time"
)

type IDB interface {
	GetCloseFunction() error
	GetConnectionPoll() *pgxpool.Pool
	GetTransaction(ctx context.Context, TxOptions pgx.TxOptions) (pgx.Tx, context.Context, error)
	CommitTransaction(ctx context.Context, transaction pgx.Tx) error
	RollbackTransaction(ctx context.Context, transaction pgx.Tx) error
}

type db struct {
	Pool             *pgxpool.Pool
	cancellationTime time.Duration
}

func NewDB(ctx context.Context, cons *pgxpool.Pool,
	CancellationTime time.Duration) IDB {
	return &db{
		Pool:             cons,
		cancellationTime: CancellationTime,
	}

}

func (d *db) GetCloseFunction() error {
	d.Pool.Close()
	return nil
}

func (d *db) GetConnectionPoll() *pgxpool.Pool {
	return d.Pool
}

func (d *db) GetTransaction(ctx context.Context, TxOptions pgx.TxOptions) (pgx.Tx, context.Context, error) {
	cancelCtx, _ := context.WithTimeout(ctx, d.cancellationTime)
	tx, err := d.Pool.BeginTx(cancelCtx, TxOptions)
	if err != nil {
		return nil, nil, err
	}

	return tx, cancelCtx, err
}

func (d *db) CommitTransaction(ctx context.Context, transaction pgx.Tx) error {
	if transaction == nil {
		return errors.New("nil transaction")
	}
	return transaction.Commit(ctx)
}

func (d *db) RollbackTransaction(ctx context.Context, transaction pgx.Tx) error {
	if transaction == nil {
		return errors.New("nil transaction")
	}
	return transaction.Rollback(ctx)
}
