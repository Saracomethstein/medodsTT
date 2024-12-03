package services

import (
	"database/sql"
	"medodsTT/internal/repositories"
)

type ServiceContainer struct {
	GenerateService *GenerateService
	RefreshService  *RefreshService
}

func NewServiceContainer(db *sql.DB) *ServiceContainer {
	generateRepo := repositories.NewGenerateRepository(db)
	refreshRepo := repositories.NewRefreshRepository(db)

	return &ServiceContainer{
		GenerateService: NewGenerateService(*generateRepo),
		RefreshService:  NewRefreshService(*refreshRepo),
	}
}
