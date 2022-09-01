package suite

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.ozon.dev/astoyakin/route256/internal/pkg/db"
	"io/ioutil"
	"strings"
	"time"
)

type DB struct {
	Database db.IDB
}

func NewDB(ctx context.Context, cons *pgxpool.Pool,
	CancellationTime time.Duration) /*db.IDB*/ *DB {
	return &DB{
		Database: db.NewDB(ctx, cons, CancellationTime),
	}
}

func (f *DB) ImportFixtures(paths ...string) error {
	ctx := context.Background()

	// run import inside a transaction
	tx, err := f.Database.GetConnectionPoll().Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// return the initial error
			_ = tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()
	// path to fixture
	for _, path := range paths {
		/* #nosec */
		file, fileErr := ioutil.ReadFile(path)
		if fileErr != nil {
			return fileErr
		}

		// split file content by ";\n"
		for _, sql := range strings.Split(string(file), ";") {
			if strings.TrimSpace(sql) != "" {
				if _, err = tx.Exec(ctx, sql); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// clean specified db tables
func (f *DB) TruncateTables(t ...string) error {
	ctx := context.Background()
	_, err := f.Database.GetConnectionPoll().Exec(ctx, fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE", strings.Join(t, ",")))

	return err
}
