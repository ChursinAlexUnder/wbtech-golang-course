package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dsn = "postgres://wbtech-L0:16530@localhost:5432/wbtech-golang-course-L0?sslmode=prefer"

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
