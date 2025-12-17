package main

import (
	"context"
	"fmt"
	"notifications-service/pkg/storage"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Notification struct {
	UserID  int
	Message string
}

var notificationsChannel = make(chan Notification, 1000)

func startNotificationWorker(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case n := <-notificationsChannel:
				saveNotification(n)

			case <-ctx.Done():
				fmt.Println("Notification worker stopped")
				return
			}
		}
	}()
}

func sendNotification(n Notification) {
	select {
	case notificationsChannel <- n:
	default:
	}
}

func saveNotification(n Notification) {
	storage.SaveNotification(n.UserID, n.Message)
	time.Sleep(time.Millisecond * 300)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	connStr := "postgres://postgres:1234@localhost:5432/notifications_db?sslmode=disable"
	err := storage.InitDB(connStr)
	if err != nil {
		panic(err)
	}

	startNotificationWorker(ctx, &wg)

	sendNotification(Notification{UserID: 1, Message: "New comment"})
	sendNotification(Notification{UserID: 2, Message: "New message"})

	for i := 0; i < 15; i++ {
		sendNotification(Notification{UserID: i, Message: fmt.Sprintf("Message #%d", i)})
	}

	fmt.Println("Server started. Press Ctrl+C to stop.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()
	wg.Wait()
	fmt.Println("Server shutdown gracefully")
}
