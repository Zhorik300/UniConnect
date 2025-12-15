package admin


import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	r.GET("/dashboard", adminDashboard)
}

func adminDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome admin",
	})
}
