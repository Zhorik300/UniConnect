package main

import (
	"log"
	"uniconnect/internal/auth"
	"uniconnect/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", auth.RegisterHandler)
			authRoutes.POST("/login", auth.LoginHandler)
		}
	}

	r.Run(":8080")
}
