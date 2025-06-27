package tag

import (
	"joosum-backend/app/user"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagUsecase TagUsecase
}

type CreateTagReq struct {
	Names []string `json:"names"`
}

// CreateTags godoc
// @Tags 태그
// @Summary 태그를 생성합니다.
// @Description 처음엔 INSERT 그 다음엔 UPDATE 가 되도록 했어요.
// @Accept  json
// @Produce  json
// @Param request body []string true "태그 생성 요청 본문"
// @Success 200 {array} []string "태그 생성이 성공적으로 이루어졌을 때 새로 생성된 태그 리스트 반환"
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않을 때 반환합니다."
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 500 {object} util.APIError "태그 생성 과정에서 오류가 발생한 경우 반환합니다."
// @Security ApiKeyAuth
// @Router /tags [post]
func (h TagHandler) CreateTags(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, util.MsgMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	var req []string
	if err := c.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	tags, err := h.tagUsecase.CreateTags(userId, req)
	if err != nil {
		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, tags.Names)
}

// GetTags godoc
// @Tags 태그
// @Summary 태그를 조회합니다.
// @Description 사용자 아이디를 통해 해당 사용자의 모든 태그를 조회합니다. 검색어를 전달하면 해당 검색어에 맞는 태그만 반환합니다.
// @Accept  json
// @Produce  json
// @Param search query string false "검색어"
// @Success 200 {array} []string "태그 조회가 성공적으로 이루어졌을 때 태그 배열 반환"
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 500 {object} util.APIError "태그 조회 과정에서 오류가 발생한 경우 반환합니다."
// @Security ApiKeyAuth
// @Router /tags [get]
func (h TagHandler) GetTags(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, util.MsgMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	search := c.Query("search")

	var (
		tags []string
		err  error
	)

	if search == "" {
		tags, err = h.tagUsecase.FindTagsByUserId(userId)
	} else {
		tags, err = h.tagUsecase.FindTagsByUserIdAndSearch(userId, search)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, []string{})
			return
		}

		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, tags)
}

// DeleteTag godoc
// @Tags 태그
// @Summary 태그를 삭제합니다.
// @Description 사용자 아이디와 태그 명을 통해 해당 태그를 삭제합니다.
// @Accept  json
// @Produce  json
// @Param tag path string true "태그 명"
// @Success 200 {array} []string "태그"
// @Failure 400 {object} util.APIError "요청이 유효하지 않을 때 반환합니다."
// @Failure 401 {object} util.APIError "Authorization 헤더가 없을 때 반환합니다."
// @Failure 500 {object} util.APIError "태그 삭제 과정에서 오류가 발생한 경우 반환합니다."
// @Security ApiKeyAuth
// @Router /tags/{tag} [delete]
func (h TagHandler) DeleteTag(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		// 401 Unauthorized
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, util.MsgMissingAuthorization)
		return
	}

	userId := currentUser.(*user.User).UserId

	tag := c.Param("tag")

	tags, err := h.tagUsecase.DeleteTag(userId, tag)
	if err != nil {
		// 500 Internal Server Error
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, tags)
}
