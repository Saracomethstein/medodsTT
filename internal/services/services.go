package services

import (
	"database/sql"
	"medodsTT/internal/repositories"
)

type ServiceContainer struct {
	GenerateService *TokenService
	RefreshService  *TokenService
}

func NewServiceContainer(db *sql.DB) *ServiceContainer {
	generateRepo := repositories.NewTokenRepository(db)
	refreshRepo := repositories.NewTokenRepository(db)

	return &ServiceContainer{
		GenerateService: NewTokenService(*generateRepo),
		RefreshService:  NewTokenService(*refreshRepo),
	}
}
