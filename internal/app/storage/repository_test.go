package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/astoyakin/route256/internal/app/storage/mocks"
	"testing"
)

type RepositoryMockStruct struct {
	pool        pgxmock.PgxPoolIface
	repo        mocks.RepositoryMock
	transaction pgx.Tx
}

func SuiteSetUp(t *testing.T) RepositoryMockStruct {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectBeginTx(pgx.TxOptions{})
	tx, err := mock.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		t.Fatal(err)
	}
	return RepositoryMockStruct{pool: mock,
		repo:        mocks.RepositoryMock{},
		transaction: tx}

}

func (rms *RepositoryMockStruct) SuiteDown(t *testing.T) {
	rms.pool.Close()
}

func Test_repository_CommitTransaction(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
	}
	tests := struct {
		name string
		args args
	}{
		name: "success",
		args: args{ctx: context.Background(),
			transaction: f.transaction},
	}

	// TODO: Add test cases.

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("CommitTransaction", tests.args.ctx, tests.args.transaction).Return(nil).Once()

		err := f.repo.CommitTransaction(tests.args.ctx, tests.args.transaction)
		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	tests = struct {
		name string
		args args
	}{
		name: "failed - transaction is rollbacked already",
		args: args{ctx: context.Background(),
			transaction: f.transaction},
	}

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("CommitTransaction", tests.args.ctx, tests.args.transaction).Return(pgx.ErrTxClosed).Once()

		err := f.repo.CommitTransaction(tests.args.ctx, tests.args.transaction)
		require.Error(t, err)
		require.Equal(t, err, pgx.ErrTxClosed)
		f.repo.AssertExpectations(t)
	})

}

func Test_repository_GetTransaction(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx       context.Context
		TxOptions pgx.TxOptions
	}
	testEx := struct {
		name string
		args args
	}{
		name: "success",
		args: args{ctx: context.Background(),
			TxOptions: pgx.TxOptions{}},
	}

	t.Run(testEx.name, func(t *testing.T) {
		f.repo.On("GetTransaction", testEx.args.ctx, testEx.args.TxOptions).Return(f.transaction, testEx.args.ctx, nil).Once()

		trans, _, err := f.repo.GetTransaction(testEx.args.ctx, testEx.args.TxOptions)
		require.NoError(t, err)
		require.NotEmpty(t, trans)
		f.repo.AssertExpectations(t)
	})

	testEx = struct {
		name string
		args args
	}{
		name: "failed - context exceeded",
		args: args{ctx: context.Background(),
			TxOptions: pgx.TxOptions{}},
	}

	t.Run(testEx.name, func(t *testing.T) {
		f.repo.On("GetTransaction", testEx.args.ctx, testEx.args.TxOptions).Return(nil, testEx.args.ctx, context.DeadlineExceeded).Once()

		trans, _, err := f.repo.GetTransaction(testEx.args.ctx, testEx.args.TxOptions)
		require.Equal(t, err, context.DeadlineExceeded)
		require.Empty(t, trans)
		f.repo.AssertExpectations(t)
	})

}

func Test_repository_RollbackTransaction(t *testing.T) {
	f := SuiteSetUp(t)
	defer f.SuiteDown(t)

	type args struct {
		ctx         context.Context
		transaction pgx.Tx
	}
	tests := struct {
		name string
		args args
	}{
		name: "success",
		args: args{ctx: context.Background(),
			transaction: f.transaction},
	}

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("RollbackTransaction", tests.args.ctx, tests.args.transaction).Return(nil).Once()

		err := f.repo.RollbackTransaction(tests.args.ctx, tests.args.transaction)
		require.NoError(t, err)
		f.repo.AssertExpectations(t)
	})

	tests = struct {
		name string
		args args
	}{
		name: "fail - transaction is commited",
		args: args{ctx: context.Background(),
			transaction: f.transaction},
	}

	t.Run(tests.name, func(t *testing.T) {
		f.repo.On("RollbackTransaction", tests.args.ctx, tests.args.transaction).Return(pgx.ErrTxClosed).Once()

		err := f.repo.RollbackTransaction(tests.args.ctx, tests.args.transaction)
		require.Error(t, err)
		require.Equal(t, err, pgx.ErrTxClosed)
		f.repo.AssertExpectations(t)
	})

}
