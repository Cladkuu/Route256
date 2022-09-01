package orderStatusStorage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/mocks"
	"testing"
)

type StatusRepositoryMockStruct struct {
	pool        pgxmock.PgxPoolIface
	repo        mocks.OrderStatusRepositoryMock
	transaction pgx.Tx
}

func SuiteSetUp(t *testing.T) StatusRepositoryMockStruct {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBeginTx(pgx.TxOptions{})
	tx, err := mock.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}

	return StatusRepositoryMockStruct{
		pool:        mock,
		repo:        mocks.OrderStatusRepositoryMock{},
		transaction: tx,
	}
}

func (crms *StatusRepositoryMockStruct) SuiteDown(t *testing.T) {
	crms.pool.Close()
}

func Test_orderStatusRepository_ValidateStatusesV2(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		id          int64
		statusTo    string
	}
	tests := struct {
		args args
	}{
		args: args{ctx: context.Background(),
			transaction: f.transaction,
			id:          1,
			statusTo:    "IN_DELIVERY"},
	}

	t.Run("success", func(t *testing.T) {
		f.repo.On("ValidateStatusesV2", tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.statusTo).Return(nil).Once()
		err := f.repo.ValidateStatusesV2(tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.statusTo)
		require.NoError(t, err)

		f.repo.AssertExpectations(t)
	})

	t.Run("fail - pass is forbidden", func(t *testing.T) {
		f.repo.On("ValidateStatusesV2", tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.statusTo).Return(sql.ErrNoRows).Once()
		err := f.repo.ValidateStatusesV2(tests.args.ctx, tests.args.transaction, tests.args.id, tests.args.statusTo)
		require.Error(t, err)
		require.Equal(t, err, sql.ErrNoRows)
		f.repo.AssertExpectations(t)
	})

}
