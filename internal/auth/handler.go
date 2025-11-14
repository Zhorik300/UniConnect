// internal/auth/handlers.go
package auth

import (
    "net/http"
    "time"
    "uniconnect/internal/database"
    "uniconnect/internal/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Email    string `json:"email,omitempty"`
}

func RegisterHandler(c *gin.Context) {
    var creds Credentials
    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

    _, err := database.DB.Exec(`
        INSERT INTO users (username, email, password) VALUES ($1, $2, $3)
    `, creds.Username, creds.Email, string(hashedPassword))

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

func LoginHandler(c *gin.Context) {
    var creds Credentials
    if err := c.ShouldBindJSON(&creds); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    err := database.DB.Get(&user, "SELECT * FROM users WHERE username=$1", creds.Username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, _ := token.SignedString(jwtKey)
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
