package link

import (
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type LinkBookHandler struct {
	linkBookUsecase LinkBookUsecase
}

// GetLinkBooks
// @Tags 링크북
// @Summary 링크북 목록 조회
// @Description `sort` 파라미터는 `created_at`, `last_saved_at`, `title` 중 하나만 허용됩니다.
// @Param sort query string false "정렬 기준" Enums(created_at,last_saved_at,title)
// @Success 200 {object} link.LinkBookListRes
// @Failure 400 {object} util.APIError "허용되지 않은 sort 값이 전달된 경우"
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 500 {object} util.APIError "링크북 목록 조회 과정에서 오류가 발생한 경우"
// @Security ApiKeyAuth
// @Router /link-books [get]
func (h LinkBookHandler) GetLinkBooks(c *gin.Context) {
	var req LinkBookListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody)
		return
	}

	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	res, err := h.linkBookUsecase.GetLinkBooks(req, userId)
	if err != nil {
		if err == util.ErrInvalidSort {
			util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody)
			return
		}
		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreateLinkBook
// @Tags 링크북
// @Summary 링크북 생성
// @Param request body link.LinkBookCreateReq true "request"
// @Success 200 {object} link.LinkBook
// @Security ApiKeyAuth
// @Router /link-books [post]
func (h LinkBookHandler) CreateLinkBook(c *gin.Context) {
	var req LinkBookCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody)
		return
	}

	err := util.Validate.Struct(req)
	if err != nil {
		fields := util.ValidatorErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    util.CodeInvalidRequestBody,
			"message": util.ErrorMessage(util.CodeInvalidRequestBody),
			"fields":  fields,
		})
		return
	}

	req.Title = strings.Trim(req.Title, " ")

	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	res, err := h.linkBookUsecase.CreateLinkBook(req, userId)
	if err != nil {
		if err == util.ErrDuplicatedTitle {
			util.SendError(c, http.StatusBadRequest, util.CodeDuplicateTitle)
			return
		}

		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// UpdateLinkBook
// @Tags 링크북
// @Summary 링크북 수정
// @Param        linkBookId   path      string  true  "LinkBookId"
// @Param request body link.LinkBookCreateReq true "request"
// @Success 200 {object} link.LinkBook
// @Security ApiKeyAuth
// @Router /link-books/{linkBookId} [put]
func (h LinkBookHandler) UpdateLinkBook(c *gin.Context) {
	var req LinkBookCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody)
		return
	}

	req.Title = strings.Trim(req.Title, " ")

	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	linkBookId := c.Param("linkBookId")

	res, err := h.linkBookUsecase.UpdateLinkBook(linkBookId, req, userId)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			util.SendError(c, http.StatusNotFound, util.CodeInvalidRequestBody)
			return
		}

		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteLinkBook
// @Tags 링크북
// @Summary 링크북 삭제
// @Description 링크북과 모든 링크들을 삭제 (기본 링크북이라면 링크들만 삭제)
// @Param        linkBookId   path      string  true  "LinkBookId"
// @Success 200 {object} link.LinkBookDeleteRes
// @Security ApiKeyAuth
// @Router /link-books/{linkBookId} [delete]
func (h LinkBookHandler) DeleteLinkBook(c *gin.Context) {
	linkBookId := c.Param("linkBookId")

	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	res, err := h.linkBookUsecase.DeleteLinkBookWithLinks(userId, linkBookId)
	if err != nil {
		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
