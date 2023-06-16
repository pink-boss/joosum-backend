package link

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/pkg/util"
	"net/http"
)

type LinkBookHandler struct {
	linkUsecase LinkUsecase
}

// CreateLinkBook
// @Tags 링크
// @Summary 링크북 생성
// @Param request body link.LinkBook true "request"
// @Success 200 {object} link.LinkBookResponse
// @Router /link-books [post]
func (h LinkBookHandler) CreateLinkBook(c *gin.Context) {
	var req LinkBook
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.APIError{Error: "Invalid request body"})
		return
	}

	res, err := h.linkUsecase.CreateLinkBook(req)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
