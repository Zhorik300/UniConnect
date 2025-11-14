package websocket

import (
	"net/http"
	"strconv"
	"uniconnect/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections = make(map[int]map[*websocket.Conn]bool) // postID -> set of connections

func CommentsWS(c *gin.Context) {
	postIDParam := c.Param("postId")
	postID, _ := strconv.Atoi(postIDParam)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	if connections[postID] == nil {
		connections[postID] = make(map[*websocket.Conn]bool)
	}
	connections[postID][conn] = true

	for {
		var msg struct {
			AuthorID int    `json:"author_id"`
			Content  string `json:"content"`
		}
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(connections[postID], conn)
			conn.Close()
			break
		}

		// Сохраняем комментарий в БД
		_, _ = database.DB.Exec(
			"INSERT INTO comments (post_id, author_id, content) VALUES ($1, $2, $3)",
			postID, msg.AuthorID, msg.Content,
		)

		// Отправляем всем подписанным на этот пост
		for c := range connections[postID] {
			c.WriteJSON(msg)
		}
	}
}
