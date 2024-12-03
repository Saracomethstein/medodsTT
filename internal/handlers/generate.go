package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

func NewAuthHandler(db *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) GenerateToken(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	accessToken, err := generateAccessToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate access token"})
	}

	refreshToken, refreshHash, err := generateRefreshToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate refresh token"})
	}

	err = saveRefreshTokenToDB(h.DB, userID, refreshHash, c.RealIP())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save refresh token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func saveRefreshTokenToDB(db *pgxpool.Pool, userID, hash, ip string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, ip_address, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(context.Background(), query, userID, hash, ip, time.Now().Add(30*24*time.Hour))
	return err
}
