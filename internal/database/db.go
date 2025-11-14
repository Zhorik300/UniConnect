// internal/database/database.go
package database

import (
    "log"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() error {
    dsn := "postgres://uniconnect_user:1234@localhost:5432/uniconnect?sslmode=disable"
    var err error
    DB, err = sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    return nil
}
