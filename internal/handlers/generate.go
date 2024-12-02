package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, "refresh token logic here")
}

func generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func generateRefreshToken() (string, string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(token)
	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return refreshToken, string(hash), nil
}

func saveRefreshTokenToDB(db *pgxpool.Pool, userID, hash, ip string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, ip_address, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(context.Background(), query, userID, hash, ip, time.Now().Add(30*24*time.Hour))
	return err
}
