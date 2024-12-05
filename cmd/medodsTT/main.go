package main

import (
	"log"
	"medodsTT/internal/handlers"
	"medodsTT/internal/repositories"
	"medodsTT/internal/services"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db := repositories.SetupDB()
	defer db.Close()
	serviceContainer := services.NewServiceContainer(db)

	generateHandler := handlers.NewGenerateHandler(serviceContainer.GenerateService)
	refreshHandler := handlers.NewRefreshHandler(serviceContainer.RefreshService)

	e.POST("/auth/token", generateHandler.GenerateToken)
	e.POST("/auth/refresh", refreshHandler.RefreshToken)

	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8000"))
}
