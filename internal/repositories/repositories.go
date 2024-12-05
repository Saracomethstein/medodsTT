package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"medodsTT/internal/models"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	retries = 5
	delay   = 3 * time.Second
)

func SetupDB() *sql.DB {
	conf := getEnv()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB_HOST, conf.DB_PORT, conf.DB_USER, conf.DB_PASSWORD, conf.DB_NAME,
	)

	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open("postgres", psqlInfo)

		if err == nil {
			err = db.Ping()

			if err == nil {
				log.Println("Successfully connected to the database.")
				return db
			}
		}

		log.Printf("Retrying to connect to the database (%d/%d): %v", i+1, retries, err)
		time.Sleep(delay)
	}

	log.Fatalf("Failed to connect to the database after %d retries: %v", retries, err)
	return nil
}

func getEnv() models.DBConnection {
	conf := new(models.DBConnection)

	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("Warning: ", err)
	}

	conf.DB_HOST = os.Getenv("DB_HOST")
	conf.DB_PORT = os.Getenv("DB_PORT")
	conf.DB_USER = os.Getenv("DB_USER")
	conf.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	conf.DB_NAME = os.Getenv("DB_NAME")

	if conf.DB_HOST == "" || conf.DB_PORT == "" || conf.DB_USER == "" || conf.DB_PASSWORD == "" || conf.DB_NAME == "" {
		log.Fatal("One or more required database environment variables are missing.")
	}

	return *conf
}
