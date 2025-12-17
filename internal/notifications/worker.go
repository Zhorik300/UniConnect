package notifications

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var NotificationsChannel = make(chan Notification, 100)

func StartWorker(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case n := <-NotificationsChannel:
				save(n)

			case <-ctx.Done():
				fmt.Println("notification worker stopped")
				return
			}
		}
	}()
}

func save(n Notification) {
	fmt.Println("Saved notification:", n.UserID, n.Message)
	time.Sleep(time.Millisecond * 200)
}
