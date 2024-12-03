package services

import "medodsTT/internal/repositories"

type RefreshService struct {
	refreshRepo repositories.RefreshRepository
}

func NewRefreshService(refreshRepo repositories.RefreshRepository) *RefreshService {
	return &RefreshService{refreshRepo: refreshRepo}
}
