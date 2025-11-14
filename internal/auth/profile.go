package auth

import (
    "net/http"
    "uniconnect/internal/database"
    "github.com/gin-gonic/gin"
)

type ProfileResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

func ProfileHandler(c *gin.Context) {
    username := c.GetString("username")
    if username == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var user ProfileResponse
    err := database.DB.Get(&user, "SELECT id, username, email, role FROM users WHERE username=$1", username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
