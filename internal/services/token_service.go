package services

import (
	"medodsTT/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt"
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
	return token.SignedString([]byte("your-secret-key"))
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
	err := s.TokenRepository.SaveRefreshToken(userID, tokenHash, ip, expiresAt)
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
