package repositories

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetupDB() (*sql.DB, error) {
	dbURL, err := loadEnv()
	if err != nil {
		return nil, err
	}

	sqldb, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	err = sqldb.Ping()
	if err != nil {
		panic(err)
	}
	return sqldb, nil
}

func loadEnv() (string, error) {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Println("No .env file found. Falling back to system environment variables.")
		return "", err
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
		return "", err
	}
	return dbURL, nil
}
