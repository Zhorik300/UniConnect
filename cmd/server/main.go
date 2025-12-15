package main

import (
	"log"
	"net/http"

	"uniconnect/internal/admin"
	"uniconnect/internal/auth"
	"uniconnect/internal/database"
	"uniconnect/internal/posts"
	"uniconnect/internal/redis"
	"uniconnect/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// DB
	if err := database.Connect(); err != nil {
        log.Fatal(err)
    }

	// Redis
	redis.Connect("localhost", "", 6379)

	// Gin
	r := gin.Default()

	// ───────────────────────────────
	// API ROOT
	// ───────────────────────────────
	api := r.Group("/api")

	// ───────────────────────────────
	// AUTH
	// ───────────────────────────────
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", auth.RegisterHandler)
		authRoutes.POST("/login", auth.LoginHandler)

		authRoutes.GET("/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Use POST /api/auth/register with JSON body: {username, password, email}",
			})
		})

		authRoutes.GET("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Use POST /api/auth/login with JSON body: {username, password}",
			})
		})

		authRoutes.GET("/profile", auth.AuthMiddleware(""), auth.ProfileHandler)
	}

	// ───────────────────────────────
// ADMIN
// ───────────────────────────────
adminRoutes := api.Group("/admin")
adminRoutes.Use(auth.AuthMiddleware("admin"))
admin.RegisterAdminRoutes(adminRoutes)


	// ───────────────────────────────
	// POSTS
	// ───────────────────────────────
	postRoutes := api.Group("/posts")
	postRoutes.Use(auth.AuthMiddleware("")) // любой авторизованный
	{
		postRoutes.POST("/", posts.CreatePostHandler)
		postRoutes.GET("/", posts.ListPostsHandler)
		postRoutes.PUT("/:id", posts.UpdatePostHandler)
		postRoutes.DELETE("/:id", posts.DeletePostHandler)
	}

	// ───────────────────────────────
	// COMMENTS
	// ───────────────────────────────
	commentRoutes := api.Group("/posts/commentary")
	commentRoutes.Use(auth.AuthMiddleware(""))
	{
		commentRoutes.POST("/:postId", posts.CreateCommentHandler)
		commentRoutes.GET("/:postId", posts.ListCommentsHandler)
	}

	// ───────────────────────────────
	// WEBSOCKETS
	// ───────────────────────────────
	r.GET("/ws/comments/:postId", websocket.CommentsWS)
	r.GET("/ws/private/:chatId", websocket.PrivateWS)

	// ───────────────────────────────

	r.Run(":8080")
}
