package posts

import (
	"net/http"
	"strconv"
	"time"
	"uniconnect/internal/database"
	"uniconnect/internal/redis"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler godoc
// @Summary Create a post
// @Description Authenticated user creates a post
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body posts.Post true "Post info"
// @Success 201 {object} posts.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts [post]
// @Security BearerAuth
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

// ListPostsHandler godoc
// @Summary List posts
// @Description Returns all posts
// @Tags Posts
// @Produce json
// @Success 200 {array} posts.Post
// @Failure 500 {object} map[string]string
// @Router /posts [get]
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

// UpdatePostHandler godoc
// @Summary Update a post
// @Description Update a post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body posts.Post true "Post data"
// @Success 200 {object} posts.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/{id} [put]
// @Security BearerAuth
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

// DeletePostHandler godoc
// @Summary Delete a post
// @Description Deletes a post by ID
// @Tags Posts
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/{id} [delete]
// @Security BearerAuth
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
