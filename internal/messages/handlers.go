package messages

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	messagesStore = make(map[string][]Message) // chatId -> messages
	msgMu         sync.Mutex
	msgSeq        int64
)

func getUserID(c *gin.Context) string {
	// Пытаемся взять user из контекста, fallback на заголовок X-User-ID (в зависимости от реализации auth)
	if v, ok := c.Get("user"); ok {
		switch u := v.(type) {
		case string:
			return u
		case int:
			return strconv.Itoa(u)
		}
	}
	if v, ok := c.Get("user_id"); ok {
		return v.(string)
	}
	if h := c.GetHeader("X-User-ID"); h != "" {
		return h
	}
	// По умолчанию пусто
	return ""
}

type sendReq struct {
	Content string `json:"content" binding:"required"`
}

func SendMessageHandler(c *gin.Context) {
	chatId := c.Param("chatId")
	var req sendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content required"})
		return
	}
	sender := getUserID(c)
	if sender == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := strconv.FormatInt(atomic.AddInt64(&msgSeq, 1), 10)
	msg := Message{
		ID:        id,
		ChatID:    chatId,
		SenderID:  sender,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	msgMu.Lock()
	messagesStore[chatId] = append(messagesStore[chatId], msg)
	msgMu.Unlock()

	c.JSON(http.StatusCreated, msg)
}

func ListMessagesHandler(c *gin.Context) {
	chatId := c.Param("chatId")
	msgMu.Lock()
	list := messagesStore[chatId]
	msgMu.Unlock()
	c.JSON(http.StatusOK, list)
}
