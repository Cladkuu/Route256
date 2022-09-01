package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/db"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/logger"
)

type repository struct {
	db db.IDB
}

func GetRepository(DB db.IDB) IRepository {
	return &repository{db: DB}
}

func (rep *repository) GetTransaction(ctx context.Context, TxOptions pgx.TxOptions) (pgx.Tx, context.Context, error) {
	tx, cancCtx, err := rep.db.GetTransaction(ctx, TxOptions)
	if err != nil {
		logger.GlobalLogger.Error(err.Error())
	}

	return tx, cancCtx, err
}

func (rep *repository) CommitTransaction(ctx context.Context, transaction pgx.Tx) error {
	if transaction == nil {
		logger.GlobalLogger.Error("nil transaction")
		return errors.New("nil transaction")
	}
	err := rep.db.CommitTransaction(ctx, transaction)
	if err != nil {
		logger.GlobalLogger.Error(err.Error())
	}
	return err
}

func (rep *repository) RollbackTransaction(ctx context.Context, transaction pgx.Tx) error {
	if transaction == nil {
		logger.GlobalLogger.Error("nil transaction")
		return errors.New("nil transaction")
	}
	err := rep.db.RollbackTransaction(ctx, transaction)
	if err != nil {
		logger.GlobalLogger.Error(err.Error())
	}
	return err
}
