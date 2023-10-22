package setting

import (
	"joosum-backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settingUsecase SettingUsecase
}

// SaveDeviceId
// @Tags 설정
// @Summary 푸시 디바이스 ID 저장
// @Param request body DeviceReq true "request"
// @Success 200 {object} db.UpdateResult
// @Security ApiKeyAuth
// @Router /settings/device [post]
func (h SettingHandler) SaveDeviceId(c *gin.Context) {
	userId := util.GetUserId(c)

	var req DeviceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.settingUsecase.SaveDeviceId(req.DeviceId, userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}

// GetNotificationAgree
// @Tags 설정
// @Summary 푸시알림 여부 조회
// @Success 200 {object} NotificationAgree
// @Security ApiKeyAuth
// @Router /settings/notification [get]
func (h SettingHandler) GetNotificationAgree(c *gin.Context) {
	userId := util.GetUserId(c)
	result, err := h.settingUsecase.GetNotificationAgree(userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}

// UpdatePushNotification
// @Tags 설정
// @Summary 푸시알림 여부 수정
// @Param request body PushNotificationReq true "request"
// @Success 200 {object} db.UpdateResult
// @Security ApiKeyAuth
// @Router /settings/notification [put]
func (h SettingHandler) UpdatePushNotification(c *gin.Context) {
	userId := util.GetUserId(c)

	var req PushNotificationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.settingUsecase.UpdatePushNotification(req, userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}
