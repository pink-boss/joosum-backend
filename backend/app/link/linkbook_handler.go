package link

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"
)

type LinkBookHandler struct {
	linkUsecase LinkUsecase
}

// GetLinkBooks
// @Tags 링크
// @Summary 링크북 목록 조회
// @Description 파라미터로 전달되는 `sort` 는 고도화 때 enum 으로 바꾸면 좋을 것 같습니다.
// @Param request query link.LinkBookListReq true "request"
// @Success 200 {object} link.LinkBookListRes
// @Security ApiKeyAuth
// @Router /link-books [get]
func (h LinkBookHandler) GetLinkBooks(c *gin.Context) {
	var req LinkBookListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: fmt.Sprintf("Invalid request parameter: %v", err.Error())})
		return
	}

	currentUser, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userId := currentUser.(*user.User).UserId

	res, err := h.linkUsecase.GetLinkBooks(req, userId)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreateLinkBook
// @Tags 링크
// @Summary 링크북 생성
// @Param request body link.LinkBookCreateReq true "request"
// @Success 200 {object} link.LinkBook
// @Security ApiKeyAuth
// @Router /link-books [post]
func (h LinkBookHandler) CreateLinkBook(c *gin.Context) {
	var req LinkBookCreateReq
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
