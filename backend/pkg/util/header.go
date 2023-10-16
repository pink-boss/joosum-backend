package util

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"net/http"
)

func GetUserId(c *gin.Context) string {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return ""
	}

	userId := currentUser.(*user.User).UserId
	return userId
}
