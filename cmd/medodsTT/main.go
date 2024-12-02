package main

import (
	"log"
	"medodsTT/internal/handlers"
	"medodsTT/internal/repositories"
	"medodsTT/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	db := repositories.SetupDB()

	serviceContainer := services.NewServiceContainer(db)

	e := echo.New()

	generateHandler := handlers.NewGenerateHandler(serviceContainer.GenerateService)
	refreshHandler := handlers.NewRefreshHandler(serviceContainer.RefreshService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/auth/token", generateHandler.GenerateToken)
	e.POST("/auth/refresh", refreshHandler.RefreshToken)

	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8000"))
}
