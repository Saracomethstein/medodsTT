package main

import (
	"fmt"
	"log"
	"medodsTT/internal/handlers"
	"medodsTT/internal/repositories"
	"medodsTT/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := repositories.SetupDB()
	if err != nil {
		fmt.Println("Can not open db connnection")
		return
	}

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
