package services

import (
	"log"
	"medodsTT/internal/models"
	"medodsTT/internal/repositories"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

type TokenService struct {
	TokenRepository repositories.TokenRepository
}

func NewTokenService(tokenRepo repositories.TokenRepository) *TokenService {
	return &TokenService{TokenRepository: tokenRepo}
}

func (s *TokenService) GenerateAccessToken(userID, ip string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"ip":  ip,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(GetJWTKey()))
}

func (s *TokenService) GenerateRefreshToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *TokenService) SaveRefreshToken(userID, tokenHash, ip string, expiresAt time.Time) error {
	err := s.TokenRepository.DeleteOldRefreshTokens(userID)
	if err != nil {
		return err
	}

	err = s.TokenRepository.SaveRefreshToken(userID, tokenHash, ip, expiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenService) GetRefreshTokenHash(userID string) string {
	token, err := s.TokenRepository.GetRefreshTokenHash(userID)
	if err != nil {
		return ""
	}
	return token
}

func (s *TokenService) GetClaimsFromJWT(request models.RefreshRequest) (string, string, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(request.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTKey()), nil
	})

	if err != nil || !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Access token is invalid or expired").SetInternal(err)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims").SetInternal(err)
	}

	ip, ok := claims["ip"].(string)
	if !ok {
		return "", "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims: missing IP").SetInternal(err)
	}
	return userID, ip, nil
}

func (s *TokenService) CompareHash(userID string, request models.RefreshRequest) error {
	refreshTokenHash := s.GetRefreshTokenHash(userID)
	err := bcrypt.CompareHashAndPassword([]byte(refreshTokenHash), []byte(request.RefreshToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token").SetInternal(err)
	}
	return nil
}

func GetJWTKey() string {
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("Warning: ", err)
	}

	key := os.Getenv("JWT_SECRET")
	if key == "" {
		log.Fatal("JWT key variables are missing.")
	}
	return key
}
