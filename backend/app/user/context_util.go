package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"joosum-backend/pkg/util"
)

// GetUserId godoc
// @Summary Context에 저장된 사용자의 ID를 반환합니다.
// @Description 미들웨어에서 설정한 "user" 정보를 조회하여 ID를 반환합니다.
// @Tags 유저
func GetUserId(c *gin.Context) string {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, util.MsgMissingAuthorization)
		return ""
	}

	userId := currentUser.(*User).UserId
	return userId
}
