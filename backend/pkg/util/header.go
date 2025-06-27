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
		SendError(c, http.StatusUnauthorized, CodeMissingAuthorization, MsgMissingAuthorization)
		return ""
	}

	userId := currentUser.(*user.User).UserId
	return userId
}
