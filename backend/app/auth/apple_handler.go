package auth

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"net/http"
)

type AppleHandler struct {
	authUsecae   AuthUsecae
	appleUsecase AppleUsecase
	userUsecase  user.UserUsecase
}

// VerifyAndIssueToken
// @Tags 로그인
// @Summary id_token 을 verify 한 후 애플로 부터 토큰 발급
// @Param request body auth.authRequest true "code 와 id_token"
// @Success 200 {object} auth.tokenResponse
// @Router /auth/apple [post]
func (h *AppleHandler) VerifyAndIssueToken(c *gin.Context) {
	reqAuth := authRequest{}
	if err := c.Bind(&reqAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding failure"})
		return
	}

	if reqAuth.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	claims, err := h.appleUsecase.VerifyAccessToken(reqAuth)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	email := claims["email"].(string)

	// 정보를 입력하고 회원가입을 했는지 확인
	user, err := h.userUsecase.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"accessToken": "", "refreshToken": ""})
		return
	}

	// 주섬토큰 발급
	accessToken, refreshToken, err := h.authUsecae.GenerateNewJWTToken([]string{"USER", "ADMIN"}, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
