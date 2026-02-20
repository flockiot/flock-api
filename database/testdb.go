package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestDB struct {
	Pool *pgxpool.Pool
}

func NewTestDB(ctx context.Context, dsn string) (*TestDB, error) {
	pool, err := Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting to test database: %w", err)
	}

	if err := Migrate(dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("running test migrations: %w", err)
	}

	return &TestDB{Pool: pool}, nil
}

func (db *TestDB) Close() {
	db.Pool.Close()
}
