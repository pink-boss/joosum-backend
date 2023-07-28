package admin

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"
)

type AdminHandler struct {
	adminUsecase AdminUsecase
}

// UpdateDefaultFolder
// @Tags 관리자
// @Summary 기본폴더 정보 변경
// @Param request body admin.LinkBookUpdateReq true "request"
// @Success 200 {object} db.LinkBook
// @Security ApiKeyAuth
// @Router /admin/link-books [put]
func (h AdminHandler) UpdateDefaultFolder(c *gin.Context) {
	var req LinkBookUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	err := util.Validate.Struct(req)
	if err != nil {
		fields := util.ValidatorErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ValidatorErrors", "fields": fields})
		return
	}

	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userId := currentUser.(*user.User).UserId
	if userId != "User-dea95e0a-6d06-4d9f-bd2e-094bcedcc792" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not an administrator"})
		return
	}

	name := currentUser.(*user.User).Name
	res, err := h.adminUsecase.UpdateDefaultFolder(req, name)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
