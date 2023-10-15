package notif

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/pkg/util"
	"net/http"
)

type NotificationHandler struct {
	notificationUsecase NotificationUsecase
}

// SaveDeviceId
// @Tags 알림
// @Summary 푸시 디바이스 ID 저장
// @Param request body DeviceReq true "request"
// @Success 200 {object} db.UpdateResult
// @Security ApiKeyAuth
// @Router /notifications/device [post]
func (h NotificationHandler) SaveDeviceId(c *gin.Context) {
	userId := util.GetUserId(c)

	var req DeviceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.notificationUsecase.SaveDeviceId(req.DeviceId, userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}
