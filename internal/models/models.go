package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), databaseURL)
}
