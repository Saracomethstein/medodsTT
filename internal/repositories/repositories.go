package repositories

import (
	"database/sql"
	"fmt"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "medods"
)

func SetupDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	sqldb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = sqldb.Ping()
	if err != nil {
		panic(err)
	}
	return sqldb
}

// func SetupDB() (*sql.DB, error) {
// 	dbURL, err := loadEnv()
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqldb, err := sql.Open("postgres", dbURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = sqldb.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return sqldb, nil
// }

// func loadEnv() (string, error) {
// 	err := godotenv.Load("/app/.env")
// 	if err != nil {
// 		log.Println("No .env file found. Falling back to system environment variables.")
// 		return "", err
// 	}

// 	dbURL := os.Getenv("DATABASE_URL")
// 	if dbURL == "" {
// 		log.Fatal("DATABASE_URL is not set")
// 		return "", err
// 	}
// 	return dbURL, nil
// }
