package services

import (
	"log"
	"medodsTT/internal/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
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
	return token.SignedString([]byte("J8sK^7z!fA0p@o3wY%M#E1Qx%Rk4U&Nv2KZ"))
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

func GetJWTKey() string {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	key := os.Getenv("JWT_SECRET")
	if key == "" {
		log.Fatal("JWT key variables are missing.")
	}
	return key
}
