package auth

import (
	"fmt"
	"joosum-backend/app/user"
	"joosum-backend/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type NaverHandler struct {
	authUsecase AuthUsecase
	naverUsecae NaverUsecae
	userUsecase user.UserUsecase
}

// @Summary Naver 액세스 토큰 검증
// @Description Naver 액세스 토큰의 유효성을 검사하고 새 JWT 토큰 쌍을 생성합니다.
// @Tags 로그인
// @Accept  json
// @Produce  json
// @Param request body NaverAccessTokenReq true "액세스 토큰 요청 본문"
// @Success 200 {object} util.TokenRes "액세스 토큰을 성공적으로 검증하고 JWT 토큰을 반환합니다."
// @Failure 400 {object} util.APIError "요청 본문이 유효하지 않거나 액세스토큰이 누락된 경우 Bad Request를 반환합니다."
// @Failure 401 {object} util.APIError "Naver 액세스 토큰이 유효하지 않거나 사용자가 존재하지 않는 경우 Unauthorized를 반환합니다."
// @Failure 500 {object} util.APIError "액세스 토큰을 검증하거나 JWT 토큰을 생성하는 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/naver [post]
func (h *NaverHandler) VerifyNaverAccessToken(c *gin.Context) {
	var req AccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, util.MsgInvalidRequestBody)
		return
	}

	accessToken := req.IdToken

	if accessToken == "" {
		util.SendError(c, http.StatusBadRequest, util.CodeInvalidRequestBody, "idToken 값이 필요합니다")
		return
	}

	email, err := h.naverUsecae.GetUserEmailByToken(accessToken)
	if err != nil {
		util.SendError(c, http.StatusUnauthorized, util.CodeMissingAuthorization, err.Error())
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil && err != mongo.ErrNoDocuments {
		util.SendError(c, http.StatusUnauthorized, util.CodeInternalServerError, fmt.Sprintf("failed to get user by email: %v", err.Error()))
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, util.TokenRes{AccessToken: "", RefreshToken: ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		util.SendError(c, http.StatusInternalServerError, util.CodeInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, util.TokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken})

}
