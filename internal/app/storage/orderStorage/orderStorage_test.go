package orderStorage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	model2 "gitlab.ozon.dev/astoyakin/route256/internal/app/model"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/mocks"
	"testing"
)

type OrderRepositoryMockStruct struct {
	pool        pgxmock.PgxPoolIface
	repo        mocks.OrderRepositoryMock
	transaction pgx.Tx
}

func SuiteSetUp(t *testing.T) OrderRepositoryMockStruct {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBeginTx(pgx.TxOptions{})
	tx, err := mock.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}

	return OrderRepositoryMockStruct{
		pool:        mock,
		repo:        mocks.OrderRepositoryMock{},
		transaction: tx,
	}
}

func (crms *OrderRepositoryMockStruct) SuiteDown(t *testing.T) {
	crms.pool.Close()
}

func Test_orderRepository_CancelOrder(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		id          int64
	}
	tests := struct {
		name string
		args args
	}{
		name: "success",
		args: args{ctx: context.Background(),
			transaction: f.transaction,
			id:          1},
	}

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("CancelOrder", tests.args.ctx, tests.args.transaction, tests.args.id).Return(nil).Once()

		err := f.repo.CancelOrder(tests.args.ctx, tests.args.transaction, tests.args.id)
		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("CancelOrder", tests.args.ctx, tests.args.transaction, tests.args.id).Return(context.DeadlineExceeded).Once()

		err := f.repo.CancelOrder(tests.args.ctx, tests.args.transaction, tests.args.id)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		f.repo.AssertExpectations(t)
	})

}

func Test_orderRepository_ChangeStatus(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		id          int64
		status      string
	}
	tests := struct {
		name string
		args args
	}{
		args: args{
			ctx:         context.Background(),
			transaction: f.transaction,
			id:          1,
			status:      "IN_DELIVERY",
		},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("ChangeStatus", tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.status).Return(nil).Once()
		err := f.repo.ChangeStatus(tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.status)
		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("ChangeStatus", tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.status).Return(context.DeadlineExceeded).Once()

		err := f.repo.ChangeStatus(tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.status)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		f.repo.AssertExpectations(t)
	})
}

func Test_orderRepository_CreateOrder(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		price       int32
		currency    string
		orderCode   string
	}
	tests := struct {
		name string
		args args
	}{
		args: args{ctx: context.Background(),
			transaction: f.transaction,
			price:       30,
			currency:    "RUB",
			orderCode:   "code"},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("CreateOrder", tests.args.ctx, tests.args.transaction, tests.args.price, tests.args.currency, tests.args.orderCode).Return(int64(1), nil).Once()
		id, err := f.repo.CreateOrder(tests.args.ctx, tests.args.transaction, tests.args.price, tests.args.currency, tests.args.orderCode)
		require.NoError(t, err)
		require.Equal(t, id, int64(1))
		f.repo.AssertExpectations(t)
	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("CreateOrder", tests.args.ctx, tests.args.transaction, tests.args.price, tests.args.currency, tests.args.orderCode).Return(int64(0), context.DeadlineExceeded).Once()
		id, err := f.repo.CreateOrder(tests.args.ctx, tests.args.transaction, tests.args.price, tests.args.currency, tests.args.orderCode)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		require.Equal(t, id, int64(0))
		f.repo.AssertExpectations(t)
	})

}

func Test_orderRepository_GetAllOrders(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx          context.Context
		transaction  pgx.Tx
		page         int32
		pageSize     int32
		sortCriteria string
	}
	tests := struct {
		name string
		args args
		want []*model2.Order
	}{
		args: args{
			ctx:          context.Background(),
			transaction:  f.transaction,
			page:         1,
			pageSize:     5,
			sortCriteria: "asc",
		},
		want: []*model2.Order{},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("GetAllOrders", tests.args.ctx, tests.args.transaction, tests.args.page, tests.args.pageSize, tests.args.sortCriteria).Return(tests.want, nil).Once()
		got, err := f.repo.GetAllOrders(tests.args.ctx, tests.args.transaction, tests.args.page, tests.args.pageSize, tests.args.sortCriteria)
		require.NoError(t, err)
		require.Equal(t, len(tests.want), len(got))
		f.repo.AssertExpectations(t)

	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("GetAllOrders", tests.args.ctx, tests.args.transaction, tests.args.page, tests.args.pageSize, tests.args.sortCriteria).Return(nil, context.DeadlineExceeded).Once()
		_, err := f.repo.GetAllOrders(tests.args.ctx, tests.args.transaction, tests.args.page, tests.args.pageSize, tests.args.sortCriteria)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		f.repo.AssertExpectations(t)
	})

}

func Test_orderRepository_GetOrderById(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		id          int64
	}
	tests := struct {
		name string

		args args
		want *model2.Order
	}{
		args: args{
			ctx:         context.Background(),
			transaction: f.transaction,
			id:          1,
		},
		want: &model2.Order{},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("GetOrderById", tests.args.ctx, tests.args.transaction, tests.args.id).Return(tests.want, nil).Once()
		got, err := f.repo.GetOrderById(tests.args.ctx, tests.args.transaction, tests.args.id)
		require.NoError(t, err)
		require.Equal(t, got, tests.want)
		f.repo.AssertExpectations(t)
	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("GetOrderById", tests.args.ctx, tests.args.transaction, tests.args.id).Return(nil, context.DeadlineExceeded).Once()
		_, err := f.repo.GetOrderById(tests.args.ctx, tests.args.transaction, tests.args.id)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		f.repo.AssertExpectations(t)
	})

}

func Test_orderRepository_ResetOrderPrice(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		id          int64
	}
	tests := struct {
		name string
		args args
	}{
		args: args{
			ctx:         context.Background(),
			transaction: f.transaction,
			id:          1,
		},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("ResetOrderPrice", tests.args.ctx, tests.args.transaction, tests.args.id).Return(nil).Once()
		err := f.repo.ResetOrderPrice(tests.args.ctx, tests.args.transaction, tests.args.id)

		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	t.Run("fail - deadline exceeded", func(t *testing.T) {
		f.repo.On("ResetOrderPrice", tests.args.ctx, tests.args.transaction, tests.args.id).Return(context.DeadlineExceeded).Once()
		err := f.repo.ResetOrderPrice(tests.args.ctx, tests.args.transaction, tests.args.id)
		require.Error(t, err)
		require.Equal(t, err, context.DeadlineExceeded)
		f.repo.AssertExpectations(t)
	})

}
