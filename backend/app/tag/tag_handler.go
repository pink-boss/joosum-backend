package tag

import (
	"joosum-backend/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagUsecae *TagUsecase
	userUsecae *user.UserUsecase
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

// CreateTag
// @Tags 태그
// @Summary 태그를 생성합니다.
// @Router /tags [post]
// @Accept  json
// @Produce  json
// @Param request body CreateTagRequest
// @Success      200  {array}   Tag
// @Failure      400  {object}  httputil.HTTPError
// @Failure      401  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
func (h TagHandler) CreateTag(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tags, err := h.tagUsecae.CreateTag(userId.(string), req.Name)
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, tags)
}

// GetTags
// @Tags 태그
// @Summary 태그를 조회합니다.
// @Router /tags [get]
// @Success      200  {array}   Tag
// @Failure      500  {object}  httputil.HTTPError
func (h TagHandler) GetTags(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	tags, err := h.tagUsecae.FindTagByUserId(userId.(string))
	if err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, tags)
}

// DeleteTag
// @Tags 태그
// @Summary 태그를 삭제합니다.
// @Param        id   path      int  true  "Account ID"
// @Router /tags/{id} [delete]
// @Success      200  {boolean} 	true
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
func (h TagHandler) DeleteTag(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		// 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	tagId := c.Param("id")

	if err := h.tagUsecae.DeleteTag(userId.(string), tagId); err != nil {
		// 500 Internal Server Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 200 OK
	c.JSON(http.StatusOK, true)
}
