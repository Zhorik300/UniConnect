package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey") // You’ll store this in .env later

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int64) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
