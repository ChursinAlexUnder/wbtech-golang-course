package database

import (
	"context"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const dsn = "postgres://wbtech-L0:16530@host.docker.internal:5432/wbtech-golang-course-L0?sslmode=disable"

// const dsn = "postgres://wbtech-L0:16530@localhost:5432/wbtech-golang-course-L0?sslmode=disable"

func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	m, err := migrate.New("file:///app/schema", dsn)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
		log.Printf("migrate.Close sourceErr=%v dbErr=%v", srcErr, dbErr)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
