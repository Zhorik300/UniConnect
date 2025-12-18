package posts

import (
	"net/http"
	"uniconnect/internal/database"

	"github.com/gin-gonic/gin"
)

// Like a post
func LikePostHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID := c.Param("id")

	_, err := database.DB.Exec(
		`INSERT INTO post_likes(post_id, user_id) VALUES($1, $2) ON CONFLICT(post_id, user_id) DO NOTHING`,
		postID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

// Unlike a post
func UnlikePostHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID := c.Param("id")

	_, err := database.DB.Exec(
		`DELETE FROM post_likes WHERE post_id=$1 AND user_id=$2`,
		postID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Like removed"})
}

// Save a post
func SavePostHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID := c.Param("id")

	_, err := database.DB.Exec(
		`INSERT INTO post_saves(post_id, user_id) VALUES($1, $2) ON CONFLICT(post_id, user_id) DO NOTHING`,
		postID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post saved"})
}

// Unsave a post
func UnsavePostHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	postID := c.Param("id")

	_, err := database.DB.Exec(
		`DELETE FROM post_saves WHERE post_id=$1 AND user_id=$2`,
		postID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post unsaved"})
}

func ListLikedPostsHandler(c *gin.Context) {
	userID := c.GetInt("user_id")

	posts := []Post{}

	err := database.DB.Select(&posts, `
        SELECT 
            p.id, 
            p.title, 
            p.content, 
            p.author_id, 
            p.created_at, 
            p.updated_at,
            COALESCE(like_counts.count, 0) AS likes_count,
            COALESCE(save_counts.count, 0) AS saved_count
        FROM posts p
        JOIN post_likes l ON l.post_id = p.id
        LEFT JOIN (
            SELECT post_id, COUNT(*) AS count
            FROM post_likes
            GROUP BY post_id
        ) AS like_counts ON like_counts.post_id = p.id
        LEFT JOIN (
            SELECT post_id, COUNT(*) AS count
            FROM post_saves
            GROUP BY post_id
        ) AS save_counts ON save_counts.post_id = p.id
        WHERE l.user_id = $1
        ORDER BY l.created_at DESC
    `, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func SearchPosts(c *gin.Context) {
	category := c.Query("category")

	var posts []Post

	// Only filter by category
	err := database.DB.Select(&posts, `
        SELECT *
        FROM posts
        WHERE ($1 = '' OR category = $1)
        ORDER BY created_at DESC
    `, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
