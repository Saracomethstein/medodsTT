package repositories

import (
	"database/sql"
	"errors"
	"time"
)

type TokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{DB: db}
}

func (r *TokenRepository) SaveRefreshToken(userID, tokenHash, ip string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, ip_address, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, token_hash) DO NOTHING
	`
	_, err := r.DB.Exec(query, userID, tokenHash, ip, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepository) GetRefreshTokenHash(userID string) (string, error) {
	var tokenHash string
	query := `SELECT token_hash FROM refresh_tokens WHERE user_id = $1 ORDER BY expires_at DESC LIMIT 1`

	err := r.DB.QueryRow(query, userID).Scan(&tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return tokenHash, nil
}
