package auth

import (
	"joosum-backend/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// https://gamz.tistory.com/11
type GoogleHandler struct {
	// DI
	// please write private
	authUsecase   AuthUsecase
	googleUsecae GoogleUsecae
	userUsecase user.UserUsecase
}

// @Summary Google 액세스 토큰 검증
// @Description Google 액세스 토큰의 유효성을 검사하고 새 JWT 토큰 쌍을 생성합니다.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body AccessTokenRequest true "액세스 토큰 요청 본문"
// @Success 200 {object} map[string]string "액세스 토큰을 성공적으로 검증하고 JWT 토큰을 반환합니다."
// @Failure 400 {object} httputil.HTTPError "요청 본문이 유효하지 않거나 액세스토큰이 누락된 경우 Bad Request를 반환합니다."
// @Failure 401 {object} httputil.HTTPError "Google 액세스 토큰이 유효하지 않거나 사용자가 존재하지 않는 경우 Unauthorized를 반환합니다."
// @Failure 500 {object} httputil.HTTPError "액세스 토큰을 검증하거나 JWT 토큰을 생성하는 과정에서 오류가 발생한 경우 Internal Server Error를 반환합니다."
// @Router /auth/google [post]
func (h *GoogleHandler) VerifyGoogleAccessToken(c *gin.Context) {
	var req AccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	accessToken := req.AccessToken

	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "accessToken is required"})
		return
	}

	if valid, err := h.googleUsecae.ValidateAccessToken(req.AccessToken); err != nil || !valid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid access token"})
		return
	}

	email, err := h.googleUsecae.GetUserEmail(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"accessToken": "", "refreshToken": ""})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken([]string{"USER", "ADMIN"}, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken})

}
