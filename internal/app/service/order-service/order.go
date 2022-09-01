package order_service

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/currencyStorage"
	orderStatusStorage2 "gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStatusStorage"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/orderStorage"
)

type orderService struct {
	orderStorage       orderStorage.IOrderStorage
	orderStatusStorage orderStatusStorage2.IOrderStatusStorage
	currencyStorage    currencyStorage.ICurrency
	commonRepository   storage.IRepository
}

func NewOrderService(OrderStorage orderStorage.IOrderStorage,
	OrderStatusStorage orderStatusStorage2.IOrderStatusStorage,
	CurrencyStorage currencyStorage.ICurrency,
	CommonRepository storage.IRepository) IOrder {
	return &orderService{
		orderStorage:       OrderStorage,
		orderStatusStorage: OrderStatusStorage,
		currencyStorage:    CurrencyStorage,
		commonRepository:   CommonRepository,
	}
}

func (os *orderService) CreateOrder(ctx context.Context, price int32, currency string, orderCode string) (int64, error) {

	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		return 0, err
	}
	defer os.commonRepository.RollbackTransaction(ctx, tx)

	/*if err = os.currencyStorage.GetCurrency(cancelCtx, tx, currency); err != nil {
		return 0, err
	}*/

	// Проверка на уникальность кода выдачи
	order, err := os.orderStorage.GetOrderByOrderCode(cancelCtx, tx, orderCode)
	/*if err != nil {
		return 0, err
	}*/
	if order != nil {
		return 0, model.NotUniqueOrderCode
	}

	orderId, err := os.orderStorage.CreateOrder(cancelCtx, tx, price,
		currency,
		orderCode)
	if err != nil {

		return 0, err
	}
	_ = os.commonRepository.CommitTransaction(ctx, tx)
	return orderId, nil
}

func (os *orderService) GetAllOrders(ctx context.Context, page, pageSize int32, sortCriteria string) ([]*model.Order, error) {
	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}

	defer os.commonRepository.RollbackTransaction(ctx, tx)

	return os.orderStorage.GetAllOrders(cancelCtx, tx, page, pageSize, sortCriteria)
}

func (os *orderService) GetOrderById(ctx context.Context, id int64) (*model.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetTransaction")
	defer span.Finish()
	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	defer os.commonRepository.RollbackTransaction(ctx, tx)

	return os.orderStorage.GetOrderById(cancelCtx, tx, id)
}

func (os *orderService) CancelOrder(ctx context.Context, id int64) error {

	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer os.commonRepository.RollbackTransaction(ctx, tx)

	if err = os.orderStatusStorage.ValidateStatusesV2(cancelCtx, tx, id, orderStatusStorage2.CancelledStatus); err != nil {
		return err
	}

	if err = os.orderStorage.CancelOrder(cancelCtx, tx, id); err != nil {
		return err
	}

	_ = os.commonRepository.CommitTransaction(ctx, tx)
	return nil
}

func (os *orderService) ChangeOrderStatus(ctx context.Context, id int64, statusTo string) error {
	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer os.commonRepository.RollbackTransaction(ctx, tx)

	if err = os.orderStatusStorage.ValidateStatusesV2(cancelCtx, tx, id, statusTo); err != nil {
		return err
	}

	err = os.orderStorage.ChangeStatus(cancelCtx, tx, id, statusTo)
	if err == nil {
		os.commonRepository.CommitTransaction(ctx, tx)
		return nil
	}

	return err
}

func (os *orderService) ResetOrderPrice(ctx context.Context, id int64) error {
	tx, cancelCtx, err := os.commonRepository.GetTransaction(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}
	defer os.commonRepository.RollbackTransaction(ctx, tx)

	order, err := os.orderStorage.GetOrderById(cancelCtx, tx, id)
	if err != nil {
		return err
	}
	if order.GetPrice() == 0 {
		return nil
	}
	err = os.orderStorage.ResetOrderPrice(cancelCtx, tx, order.GetId())
	if err == nil {
		err = os.commonRepository.CommitTransaction(ctx, tx)
		return nil
	}
	return err
}
