package orderStatusStorage

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"sync"
)

const (
	NewStatus         = "NEW"
	inPackagingStatus = "IN_PACKAGING"
	inDeliveryStatus  = "IN_DELIVERY"
	receivedStatus    = "RECEIVED"
	CancelledStatus   = "CANCELLED"
)

type orderStatusRepository struct {
	mutex *sync.Mutex

	db *pgxpool.Pool
}

const (
	statusMappingTable               = "public.status_mapping as sm"
	statusMappingTableColumnsToFetch = "sm.id_to"
	joinQuery                        = "public.status as st on st.id=sm.id_from"
)

var (
	psqlStatementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func GetOrderStatusRepository(DB *pgxpool.Pool) IOrderStatusStorage {
	orderRep := orderStatusRepository{
		mutex: &sync.Mutex{},

		db: DB,
	}
	return &orderRep
}

func (os *orderStatusRepository) FindStatuses(statusFrom, statusTo string) (bool, error) {
	os.mutex.Lock()
	defer os.mutex.Unlock()

	/*statusFromArray, ok := os.storage[statusFrom]
	if !ok {
		return false, model.UnknownStatus
	}
	for _, val := range statusFromArray {
		if val == statusTo {
			return true, nil
		}
	}*/
	return false, nil

}

func (os *orderStatusRepository) ValidateStatuses(ctx context.Context, transaction pgx.Tx, statusFrom, statusTo string) error {
	query, args, err := psqlStatementBuilder.Select(statusMappingTableColumnsToFetch).
		From(statusMappingTable).
		Join(joinQuery).
		Where(squirrel.And{
			squirrel.Eq{"id_from": psqlStatementBuilder.Select("order_id").From("public.order").Where(squirrel.Eq{"id": "id"})},
			squirrel.Eq{"id_to": psqlStatementBuilder.Select("id").From("public.status").Where(squirrel.Eq{"status": statusTo})},
		}).ToSql()
	if err != nil {
		return err
	}

	var id_to int64
	if err = pgxscan.Get(ctx, transaction, &id_to, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			err = sql.ErrNoRows
		}
		return err
	}
	return nil

}

func (os *orderStatusRepository) ValidateStatusesV2(ctx context.Context, transaction pgx.Tx, id int64, statusTo string) error {
	const query = "SELECT sm.id_to FROM public.status_mapping as sm" +
		" WHERE (id_from = (select status_id from public.orders where id=$1) AND id_to = (select id from public.status where status=$2))"

	span, ctxSpan := opentracing.StartSpanFromContext(ctx, "GetCurrency")
	defer span.Finish()

	var id_to int64
	if err := pgxscan.Get(ctxSpan, transaction, &id_to, query, id, statusTo); err != nil {
		if pgxscan.NotFound(err) {
			err = sql.ErrNoRows
			// TODO err=model.ForbiddenToPass
		}
		return err
	}
	return nil
}
