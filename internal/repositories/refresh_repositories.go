package repositories

import "database/sql"

type RefreshRepository struct {
	db *sql.DB
}

func NewRefreshRepository(db *sql.DB) *RefreshRepository {
	return &RefreshRepository{db: db}
}
