package notif

import (
	"net/http"
	. "strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
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
	userId := user.GetUserId(c)
	page, err := ParseInt(c.Query("page"), 10, 64)
	result, err := h.notificationUsecase.Notifications(userId, page)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, result)
}

// ReadNotification
// @Tags 알림
// @Summary 알림 읽음처리
// @Param notificationId path string true "알림 ID"
// @Success 204 {object} nil
// @Security ApiKeyAuth
// @Router /notifications/{notificationId} [put]
func (h NotificationHandler) ReadNotification(c *gin.Context) {
	notificationId := c.Param("notificationId")

	err := h.notificationUsecase.ReadNotification(notificationId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
