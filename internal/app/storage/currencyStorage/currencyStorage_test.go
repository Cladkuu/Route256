package currencyStorage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/mocks"
	"testing"
)

type CurrencyRepositoryMockStruct struct {
	pool        pgxmock.PgxPoolIface
	repo        mocks.CurrencyRepositoryMock
	transaction pgx.Tx
}

func SuiteSetUp(t *testing.T) CurrencyRepositoryMockStruct {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBeginTx(pgx.TxOptions{})
	tx, err := mock.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}

	return CurrencyRepositoryMockStruct{
		pool:        mock,
		repo:        mocks.CurrencyRepositoryMock{},
		transaction: tx,
	}
}

func (crms *CurrencyRepositoryMockStruct) SuiteDown(t *testing.T) {
	crms.pool.Close()
}

func Test_currencyRepository_GetCurrency(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
		currency    string
	}
	tests := struct {
		name string
		args args
	}{
		name: "success",
		args: args{
			ctx:         context.Background(),
			transaction: f.transaction,
			currency:    "RUB",
		},
	}

	t.Run(tests.name, func(t *testing.T) {
		// SELECT iso_code FROM public.currency WHERE iso_code = $1
		f.pool.ExpectQuery(`SELECT iso_code FROM public.currency WHERE iso_code = RUB`)
		f.repo.On("GetCurrency", tests.args.ctx, tests.args.transaction, tests.args.currency).Return(nil).Once()

		err := f.repo.GetCurrency(tests.args.ctx, tests.args.transaction, tests.args.currency)
		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	tests = struct {
		name string
		args args
	}{
		name: "fail - not such currency",
		args: args{
			ctx:         context.Background(),
			transaction: f.transaction,
			currency:    "RUBBB",
		},
	}

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("GetCurrency", tests.args.ctx, tests.args.transaction, tests.args.currency).Return(sql.ErrNoRows).Once()

		err := f.repo.GetCurrency(tests.args.ctx, tests.args.transaction, tests.args.currency)
		require.Error(t, err)
		require.Equal(t, err, sql.ErrNoRows)
		f.repo.AssertExpectations(t)
	})

}
