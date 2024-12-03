package services

import "medodsTT/internal/repositories"

type GenerateService struct {
	generateRepo repositories.GenerateRepository
}

func NewGenerateService(generateRepo repositories.GenerateRepository) *GenerateService {
	return &GenerateService{generateRepo: generateRepo}
}

func GenerateAccessToken(userID, ip string) (string, error) {
	return "", nil
}

func GenerateRefreshToken() (string, string, error) {
	return "", "", nil
}
