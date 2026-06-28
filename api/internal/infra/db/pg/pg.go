package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pkgerrors "github.com/pkg/errors"
)

const connectTimeout = 10 * time.Second

// MustConnect opens and pings the pool, panicking on failure. Call only from main.
func MustConnect(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	pool, err := New(ctx, databaseURL)
	if err != nil {
		panic(fmt.Sprintf("connect postgres: %+v", err))
	}
	return pool
}

// New opens a PostgreSQL connection pool and verifies connectivity with a ping.
func New(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, pkgerrors.WithStack(err)
	}

	return pool, nil
}
