package main

import (
	"log"
	"net/http"

	"uniconnect/internal/auth"
	"uniconnect/internal/database"
	"uniconnect/internal/groups"
	"uniconnect/internal/messages"
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

	// setup (создание админа по токену)
	setup := api.Group("/setup")
	{
		setup.POST("/create-admin", auth.CreateAdminSetupHandler)
	}

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
	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Welcome admin"})
		})
	}

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
		postRoutes.GET("/search", posts.SearchPosts)

		// Likes
		postRoutes.POST("/:id/like", posts.LikePostHandler)
		postRoutes.DELETE("/:id/like", posts.UnlikePostHandler)
		postRoutes.GET("/liked", posts.ListLikedPostsHandler)

		// Saves
		postRoutes.POST("/:id/save", posts.SavePostHandler)
		postRoutes.DELETE("/:id/save", posts.UnsavePostHandler)
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
	// MESSAGES (личные сообщения)
	// ───────────────────────────────
	messageRoutes := api.Group("/messages")
	messageRoutes.Use(auth.AuthMiddleware(""))
	{
		messageRoutes.POST("/:chatId", messages.SendMessageHandler)
		messageRoutes.GET("/:chatId", messages.ListMessagesHandler)
	}

	// ───────────────────────────────
	// GROUPS (группы/чат-группы и заявки)
	// ───────────────────────────────
	groupRoutes := api.Group("/groups")
	groupRoutes.Use(auth.AuthMiddleware(""))
	{
		groupRoutes.POST("/:groupId/join", groups.RequestJoinHandler) // студент — запрос на вступление
		groupRoutes.GET("/", groups.ListGroupsHandler)                // список групп (доступно всем авторизованным)
	}

	// Admin: управление группами и заявками
	adminRoutes.POST("/groups", groups.CreateGroupHandler)                             // создать группу
	adminRoutes.GET("/groups", groups.ListGroupsHandler)                               // список всех групп
	adminRoutes.GET("/groups/requests", groups.ListJoinRequestsHandler)                // список заявок
	adminRoutes.POST("/groups/requests/:id/approve", groups.ApproveJoinRequestHandler) // подтвердить заявку
	adminRoutes.DELETE("/groups/requests/:id", groups.RemoveJoinRequestHandler)        // удалить/отклонить заявку

	// ───────────────────────────────
	// WEBSOCKETS
	// ───────────────────────────────
	r.GET("/ws/comments/:postId", websocket.CommentsWS)
	r.GET("/ws/private/:chatId", websocket.PrivateWS)

	// ───────────────────────────────

	r.Run(":8080")
}
