package page

import (
	"joosum-backend/app/link"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	LinkBookUsecase link.LinkBookUsecase
	LinkUsecase     link.LinkUsecase
}

type MainPageRes struct {
	LinkBookList []link.LinkBookRes `json:"link_book_list"`
	LinkList     []link.Link        `json:"link_list"`
}

// GetMainPage godoc
// @Tags 페이지
// @Summary 메인 페이지의 데이터를 가져옵니다.
// @Description 메인 페이지의 데이터를 가져옵니다.
// @Accept  json
// @Produce  json
// @Success 200 {object} MainPageRes "메인 페이지 데이터를 성공적으로 가져왔을 때 반환합니다."
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 500 {object} util.APIError "메인 페이지 데이터를 가져오는 과정에서 오류가 발생한 경우 반환합니다."
// @Router /page/main [get]
func (h PageHandler) GetMainPage(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	/*
		linkBookList, err := h.LinkBookUsecase.FindAllLinkBooksByUserId(userId.(string))
		if err != nil {
			// 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	*/

	_, err := h.LinkUsecase.FindAllLinksByUserId(userId.(string))
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var linkBookList = []link.LinkBookRes{
		{
			Title:           "링크북1",
			BackgroundColor: "#ffffff",
			TitleColor:      "#000000",
			Illustration:    "1",
		},
	}

	var linkList = []link.Link{
		{
			LinkId:     "1",
			URL:        "https://www.naver.com",
			Title:      "네이버",
			LinkBookId: "1",
		},
	}

	res := MainPageRes{
		LinkBookList: linkBookList,
		LinkList:     linkList,
	}

	// 200 OK
	c.JSON(http.StatusOK, res)
}
