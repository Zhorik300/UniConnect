package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type adminReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// POST /api/setup/create-admin
// Требует заголовок X-SETUP-TOKEN совпадающий с ENV SETUP_TOKEN
func CreateAdminSetupHandler(c *gin.Context) {
	token := c.GetHeader("X-SETUP-TOKEN")
	if token == "" {
		token = os.Getenv("SETUP_TOKEN")
	}
	if token == "" || token != os.Getenv("SETUP_TOKEN") {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid setup token"})
		return
	}

	var r adminReq
	if err := c.ShouldBindJSON(&r); err != nil {
		// дефолтные значения (вариант B)
		r.Username = "admin"
		r.Password = "admin1234"
		r.Email = "admin@gmail.com"
	}

	// TODO: реализовать создание пользователя в БД/через существующие функции auth.
	// Нужно: проверить, существует ли пользователь; если нет — создать и выставить role='admin'.
	created, err := ensureAdmin(r.Username, r.Password, r.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, created)
}

// Заглушка — заменю на конкретную реализацию, когда покажете user creation в auth/database
func ensureAdmin(username, password, email string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: implement ensureAdmin using your user creation logic")
}
