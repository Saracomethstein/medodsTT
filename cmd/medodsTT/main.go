package main

import (
	"log"
	"medodsTT/internal/handlers"
	"medodsTT/internal/models"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("No .env file found. Falling back to system environment variables.")
		return
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
		return
	}

	db, err := models.ConnectDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	authHandler := handlers.NewAuthHandler(db)
	e.POST("/auth/token", authHandler.GenerateToken)
	e.POST("/auth/refresh", authHandler.RefreshToken)

	log.Println("Server is running on port 8080")
	log.Fatal(e.Start(":8000"))
}
