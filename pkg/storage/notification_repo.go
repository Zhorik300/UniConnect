package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDB(connString string) error {
	var err error
	db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}
	return db.Ping(context.Background())
}

func SaveNotification(userID int, message string) {
	if db == nil {
		fmt.Println("DB not initialized")
		return
	}

	_, err := db.Exec(context.Background(),
		"INSERT INTO notifications (user_id, message) VALUES ($1, $2)", userID, message)
	if err != nil {
		fmt.Println("Error saving notification:", err)
	} else {
		fmt.Println("Saved notification to DB:", userID, message)
	}
}
