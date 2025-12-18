package groups

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type Group struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Members []string  `json:"members"`
	Created time.Time `json:"created_at"`
}

type JoinRequest struct {
	ID      string    `json:"id"`
	GroupID string    `json:"group_id"`
	UserID  string    `json:"user_id"`
	Created time.Time `json:"created_at"`
}

var (
	groupMu     sync.Mutex
	groupsStore = make(map[string]*Group)
	groupSeq    int64

	reqMu    sync.Mutex
	requests = make(map[string]*JoinRequest)
	reqSeq   int64
)

func getUserID(c *gin.Context) string {
	if v, ok := c.Get("user"); ok {
		switch u := v.(type) {
		case string:
			return u
		case int:
			return strconv.Itoa(u)
		}
	}
	if h := c.GetHeader("X-User-ID"); h != "" {
		return h
	}
	return ""
}

// CreateGroupHandler — админ создаёт группу
type createGroupReq struct {
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members"`
}

func CreateGroupHandler(c *gin.Context) {
	var req createGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	id := strconv.FormatInt(atomic.AddInt64(&groupSeq, 1), 10)
	g := &Group{
		ID:      id,
		Name:    req.Name,
		Members: req.Members,
		Created: time.Now(),
	}
	groupMu.Lock()
	groupsStore[id] = g
	groupMu.Unlock()
	c.JSON(http.StatusCreated, g)
}

// ListGroupsHandler — вернуть все группы (админ) или группы в которых состоит пользователь
func ListGroupsHandler(c *gin.Context) {
	user := getUserID(c)
	groupMu.Lock()
	defer groupMu.Unlock()
	out := []*Group{}
	for _, g := range groupsStore {
		if user == "" {
			out = append(out, g)
			continue
		}
		// если пользователь задан — показать только его группы
		found := false
		for _, m := range g.Members {
			if m == user {
				found = true
				break
			}
		}
		if found {
			out = append(out, g)
		}
	}
	c.JSON(http.StatusOK, out)
}

// RequestJoinHandler — студент отправляет заявку на вступление в группу
func RequestJoinHandler(c *gin.Context) {
	user := getUserID(c)
	if user == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	groupId := c.Param("groupId")
	groupMu.Lock()
	g, ok := groupsStore[groupId]
	groupMu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	// если уже в группе
	for _, m := range g.Members {
		if m == user {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already in group"})
			return
		}
	}
	// создаём заявку
	id := strconv.FormatInt(atomic.AddInt64(&reqSeq, 1), 10)
	r := &JoinRequest{
		ID:      id,
		GroupID: groupId,
		UserID:  user,
		Created: time.Now(),
	}
	reqMu.Lock()
	requests[id] = r
	reqMu.Unlock()
	c.JSON(http.StatusCreated, r)
}

// ListJoinRequestsHandler — админ видит все заявки
func ListJoinRequestsHandler(c *gin.Context) {
	reqMu.Lock()
	defer reqMu.Unlock()
	out := []*JoinRequest{}
	for _, r := range requests {
		out = append(out, r)
	}
	c.JSON(http.StatusOK, out)
}

// ApproveJoinRequestHandler — админ подтверждает заявку
func ApproveJoinRequestHandler(c *gin.Context) {
	id := c.Param("id")
	reqMu.Lock()
	r, ok := requests[id]
	if !ok {
		reqMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	delete(requests, id)
	reqMu.Unlock()

	// добавить пользователя в группу
	groupMu.Lock()
	g, ok := groupsStore[r.GroupID]
	if !ok {
		groupMu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	g.Members = append(g.Members, r.UserID)
	groupMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "approved", "request": r})
}

// RemoveJoinRequestHandler — админ отклоняет/удаляет заявку
func RemoveJoinRequestHandler(c *gin.Context) {
	id := c.Param("id")
	reqMu.Lock()
	_, ok := requests[id]
	if ok {
		delete(requests, id)
	}
	reqMu.Unlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}
