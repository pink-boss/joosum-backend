package notif

import (
	"joosum-backend/pkg/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationUsecase NotificationUsecase
}

// Notifications
// @Tags 알림
// @Summary 알림 목록 조회
// @Param page query int false "페이지"
// @Success 200 {object} notificationResDocs
// @Security ApiKeyAuth
// @Router /notifications [get]
func (h NotificationHandler) Notifications(c *gin.Context) {
	userId := util.GetUserId(c)
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	result, err := h.notificationUsecase.Notifications(userId, page)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}
