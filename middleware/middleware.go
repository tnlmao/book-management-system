package middleware

import (
	"book-management-system/config/database"
	"book-management-system/config/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DbMiddleware(c *gin.Context) {
	err := database.ConnectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Database Connection Error")
	}
	c.Next()
}
func RedisMiddleware(c *gin.Context) {
	err := redis.ConnectToRedis()
	if err != nil {

		c.JSON(http.StatusInternalServerError, "Redis Connection Error")
		return
	}
	c.Next()
}
