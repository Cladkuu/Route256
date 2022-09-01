package currencyStorage

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	"sync"
)

type currencyRepository struct {
	mutex *sync.Mutex
	db    *pgxpool.Pool
}

var (
	psqlStatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

const (
	currencyTable   = "public.currency"
	currencyColumns = "iso_code"
)

func GetCurrencyStorage(DB *pgxpool.Pool) ICurrency {
	currencyRep := &currencyRepository{
		mutex: &sync.Mutex{},

		db: DB,
	}

	return currencyRep
}

func (cs *currencyRepository) GetCurrency(ctx context.Context, transaction pgx.Tx, currency string) error {
	query, args, err := psqlStatementBuilder.Select(currencyColumns).
		From(currencyTable).
		Where(sq.Eq{currencyColumns: currency}).
		ToSql()
	if err != nil {
		return err
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "GetCurrency")
	defer span.Finish()
	span.LogFields(log.String("currency", currency))

	var cur string
	if err = pgxscan.Get(ctx, transaction, &cur, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			err = model.UnknownCurrency
			// TODO err=model.UnknownCurrency
		}
		return err
	}

	return nil
}
