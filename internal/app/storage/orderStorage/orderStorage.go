package orderStorage

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	model2 "gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	"strings"
	"sync"
)

type orderRepository struct {
	mutex *sync.Mutex
	db    *pgxpool.Pool
}

const (
	orderTable              = "public.orders as ord"
	orderTableForInsert     = "public.orders"
	orderColumns            = "st.status,ord.price,ord.currency"
	orderJoinString         = "public.status as st ON ord.status_id=st.id"
	orderTableSortColumn    = "ord.id"
	orderTableInsertColumns = "price,currency,status_id,order_code"
)

var (
	psqlStatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

func GetOrderStorage(DB *pgxpool.Pool) IOrderStorage {
	orderRep := orderRepository{
		mutex: &sync.Mutex{},
		db:    DB,
	}

	return &orderRep
}

func (os *orderRepository) GetAllOrders(ctx context.Context, transaction pgx.Tx, page, pageSize int32, sortCriteria string) ([]*model2.Order, error) {

	query, args, err := psqlStatementBuilder.Select(orderColumns).
		From(orderTable).
		Join(orderJoinString).
		OrderBy(orderTableSortColumn + " " + strings.ToLower(sortCriteria)).
		Limit(uint64(pageSize)).
		Offset(uint64(page*pageSize - pageSize)).
		ToSql()
	if err != nil {
		return nil, err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "GetAllOrders")
	defer span.Finish()

	var orders []*model2.Order
	if err = pgxscan.Select(ctxSpan, transaction, &orders, query, args...); err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *orderRepository) GetOrderById(ctx context.Context, transaction pgx.Tx, id int64) (*model2.Order, error) {

	query, args, err := psqlStatementBuilder.Select("ord.id," + orderColumns).
		From(orderTable).
		Where(sq.Eq{"ord.id": id}).
		Join(orderJoinString).
		ToSql()
	if err != nil {
		return nil, err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "GetOrderById")
	defer span.Finish()

	var order model2.Order
	if err = pgxscan.Get(ctxSpan, transaction, &order, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			err = sql.ErrNoRows
		}
		return nil, err
	}

	return &order, nil
}

func (os *orderRepository) CreateOrder(ctx context.Context, transaction pgx.Tx, price int32, currency, orderCode string) (int64, error) {

	query, args, err := psqlStatementBuilder.Insert(orderTableForInsert).
		Columns(orderTableInsertColumns).
		Values(price, currency, 1, orderCode).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "CreateOrder")
	defer span.Finish()

	// TODO проверить отмену контекста
	var orderId int64
	if err = transaction.QueryRow(ctxSpan, query, args...).Scan(&orderId); err != nil {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		return 0, err
	}

	return orderId, nil
}

func (os *orderRepository) ChangeStatus(ctx context.Context, transaction pgx.Tx, id int64, status string) error {
	query, args, err := psqlStatementBuilder.Update(orderTableForInsert).
		Set("status_id", sq.Expr("(SELECT id FROM public.status WHERE status=$2)", status)).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "ChangeStatus")
	defer span.Finish()
	if _, err = transaction.Exec(ctxSpan, query, args[1], args[0]); err != nil {
		return err
	}
	return nil

}

func (os *orderRepository) CancelOrder(ctx context.Context, transaction pgx.Tx, id int64) error {

	query, args, err := psqlStatementBuilder.Update(orderTableForInsert).
		Set("status_id", 5).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "CancelOrder")
	defer span.Finish()
	if _, err = transaction.Exec(ctxSpan, query, args...); err != nil {
		return err
	}
	return nil
}

func (os *orderRepository) ResetOrderPrice(ctx context.Context, transaction pgx.Tx, id int64) error {

	query, args, err := psqlStatementBuilder.Update(orderTableForInsert).
		Set("price", 0).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "ResetOrderPrice")
	defer span.Finish()
	if _, err = transaction.Exec(ctxSpan, query, args...); err != nil {
		return err
	}
	return nil

}

func (os *orderRepository) GetOrderByOrderCode(ctx context.Context, transaction pgx.Tx, orderCode string) (*model2.Order, error) {
	query, args, err := psqlStatementBuilder.Select("ord.id").
		From(orderTable).
		Where(sq.Eq{"ord.order_code": orderCode}).
		ToSql()
	if err != nil {
		return nil, err
	}

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "GetOrderByOrderCode")
	defer span.Finish()

	var order model2.Order
	if err = pgxscan.Get(ctxSpan, transaction, &order, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			err = sql.ErrNoRows
		}
		return nil, err
	}

	return &order, nil
}
