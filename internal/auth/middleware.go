package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token"})
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        role := claims["role"].(string)
        if requiredRole != "" && role != requiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            return
        }
        c.Set("user_id", int(claims["user_id"].(float64)))
        c.Set("username", claims["username"])
        c.Set("role", role)
        c.Next()
    }
}
