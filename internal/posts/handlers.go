package posts

import (
	"net/http"
	"strconv"
	"time"
	"uniconnect/internal/database"
	"uniconnect/internal/redis"

	"github.com/gin-gonic/gin"
)

func CreatePostHandler(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := c.GetString("username")
	var authorID int
	err := database.DB.Get(&authorID, "SELECT id FROM users WHERE username=$1", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "author not found"})
		return
	}

	post.AuthorID = authorID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	res, err := database.DB.NamedExec(`
    INSERT INTO posts (title, content, category, author_id, created_at, updated_at)
    VALUES (:title, :content, :category, :author_id, :created_at, :updated_at)
`, &post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	post.ID = int(id)

	redis.Rdb.LPush(redis.Ctx, "notifications", username+" created a post: "+post.Title)

	c.JSON(http.StatusOK, post)
}

// ----------------- LIST -----------------
func ListPostsHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var posts []Post
	err := database.DB.Select(&posts, "SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// ----------------- UPDATE -----------------
func UpdatePostHandler(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")

	// Проверяем авторство
	var authorID int
	err := database.DB.Get(&authorID, "SELECT author_id FROM posts WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Получаем ID пользователя
	var userID int
	err = database.DB.Get(&userID, "SELECT id FROM users WHERE username=$1", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	// Проверка: либо автор, либо админ
	role := c.GetString("role")
	if userID != authorID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only update your own posts"})
		return
	}

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.UpdatedAt = time.Now()
	_, err = database.DB.NamedExec(`UPDATE posts SET title=:title, content=:content, updated_at=:updated_at WHERE id=:id`,
		map[string]interface{}{
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": post.UpdatedAt,
			"id":         id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post updated"})
}

// ----------------- DELETE -----------------
func DeletePostHandler(c *gin.Context) {
	id := c.Param("id")
	username := c.GetString("username")

	// Проверяем авторство
	var authorID int
	err := database.DB.Get(&authorID, "SELECT author_id FROM posts WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	// Получаем ID пользователя
	var userID int
	err = database.DB.Get(&userID, "SELECT id FROM users WHERE username=$1", username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	role := c.GetString("role")
	if userID != authorID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own posts"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM posts WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted"})
}
