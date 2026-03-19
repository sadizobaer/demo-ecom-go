package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDatabase() (*pgxpool.Pool, error) {
	connStr := "postgres://localhost:5432/ecom_db"

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println("Database connection successful")

	return pool, nil
}
