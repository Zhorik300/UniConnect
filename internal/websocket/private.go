package websocket

import (
	"uniconnect/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var connectionsPrivate = make(map[string]map[*websocket.Conn]bool) // chatID -> set

func PrivateWS(c *gin.Context) {
	chatID := c.Param("chatId") // например "6_7"

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	if connectionsPrivate[chatID] == nil {
		connectionsPrivate[chatID] = make(map[*websocket.Conn]bool)
	}
	connectionsPrivate[chatID][conn] = true

	for {
		var msg struct {
			SenderID   int    `json:"sender_id"`
			ReceiverID int    `json:"receiver_id"`
			Content    string `json:"content"`
		}
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(connectionsPrivate[chatID], conn)
			conn.Close()
			break
		}

		// Сохраняем сообщение в БД
		_, _ = database.DB.Exec(
			"INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3)",
			msg.SenderID, msg.ReceiverID, msg.Content,
		)

		// Отправляем всем в этом чате
		for c := range connectionsPrivate[chatID] {
			c.WriteJSON(msg)
		}
	}
}
