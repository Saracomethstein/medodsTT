package repositories

import "database/sql"

type GenerateRepository struct {
	db *sql.DB
}

func NewGenerateRepository(db *sql.DB) *GenerateRepository {
	return &GenerateRepository{db: db}
}
