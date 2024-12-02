package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.POST("/auth/token", handlers.GenerateToken())
	// e.POST("/auth/refresh", handlers.RefreshToken())

	e.Logger.Fatal(e.Start(":8080"))
}
