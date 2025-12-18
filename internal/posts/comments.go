package posts

import (
	"net/http"
	"strconv"
	"time"

	"uniconnect/internal/database"

	"github.com/gin-gonic/gin"
)

// ==========================
//
//	COMMENT MODEL
//
// ==========================
type Comment struct {
	ID        int       `db:"id" json:"id"`
	PostID    int       `db:"post_id" json:"post_id"`
	AuthorID  int       `db:"author_id" json:"author_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ==========================
//       COMMENT HANDLERS
// ==========================

// CreateCommentHandler creates a comment for a post
// @Summary Create a comment
// @Description Creates a new comment under a given post
// @Tags posts
// @Param postId path int true "Post ID"
// @Param comment body Comment true "Comment content"
// @Success 201 {object} Comment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/commentary/{postId} [post]
// @Security ApiKeyAuth
func CreateCommentHandler(c *gin.Context) {
	db := database.DB
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	userID := c.GetInt("user_id") // получаем из JWT middleware

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := Comment{
		PostID:    postID,
		AuthorID:  userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO comments (post_id, author_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = db.QueryRow(query, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt).Scan(&comment.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// ListCommentsHandler lists all comments for a post
// @Summary List comments
// @Description Returns all comments under a given post, ordered by creation time
// @Tags posts
// @Param postId path int true "Post ID"
// @Success 200 {array} Comment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /posts/commentary/{postId} [get]
// @Security ApiKeyAuth
func ListCommentsHandler(c *gin.Context) {
	db := database.DB
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post ID"})
		return
	}

	comments := []Comment{}
	err = db.Select(&comments, "SELECT * FROM comments WHERE post_id=$1 ORDER BY created_at ASC", postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}
