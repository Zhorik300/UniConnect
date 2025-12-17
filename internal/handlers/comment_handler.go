package handlers

import (
	"net/http"

	"notifications-service/internal/notifications"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	notifications.Send(notifications.Notification{
		UserID:  1,
		Message: "New comment created",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("comment created"))
}
