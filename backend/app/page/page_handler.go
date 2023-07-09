package page

import (
	"joosum-backend/app/link"
	"joosum-backend/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	linkUsecase     link.LinkUsecase
	linkBookUsecase link.LinkBookUsecase
}

type MainPageRes struct {
	LinkBookList []link.LinkBookRes
	LinkList     []*link.Link
}

// GetMainPage godoc
// @Summary 메인 페이지
// @Description 로그인한 사용자의 링크북 목록과 최근 9개의 링크를 반환합니다.
// @Tags 페이지
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} MainPageRes "메인 페이지를 성공적으로 불러오면 링크북 목록과 최근 9개의 링크를 반환합니다."
// @Failure 401 {object} util.APIError "Authorization 헤더가 없는 경우 Unauthorized를 반환합니다."
// @Failure 500 {object} util.APIError "링크북 목록 또는 링크 목록을 불러오는 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."

func (h PageHandler) GetMainPage(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	userId := currentUser.(*user.User).UserId

	linkBookList, err := h.linkBookUsecase.GetLinkbooksForMainPage(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	linkList, err := h.linkUsecase.Get9LinksByUserId(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MainPageRes{
		LinkBookList: linkBookList,
		LinkList:     linkList,
	})

}
