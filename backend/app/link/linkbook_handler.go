package link

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"
)

type LinkBookHandler struct {
	linkUsecase LinkUsecase
}

// CreateLinkBook
// @Tags 링크
// @Summary 링크북 생성
// @Param request body link.LinkBookReq true "request"
// @Success 200 {object} link.LinkBookRes
// @Security ApiKeyAuth
// @Router /link-books [post]
func (h LinkBookHandler) CreateLinkBook(c *gin.Context) {
	var req LinkBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	currentUser, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userId := currentUser.(*user.User).UserId

	res, err := h.linkUsecase.CreateLinkBook(req, userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
