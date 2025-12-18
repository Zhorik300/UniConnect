package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// дефолт для docker-compose
		dsn = "postgres://postgres:postgres@uniconnect-db:5432/uniconnect?sslmode=disable"
	}

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
