package auth

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"joosum-backend/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppleHandler struct {
	authUsecase  AuthUsecase
	appleUsecase AppleUsecase
	userUsecase  user.UserUsecase
}

// VerifyAndIssueToken
// @Tags 로그인
// @Summary 애플로그인 후 받은 id token 을 verify 한 후 주섬 JWT 토큰 발급
// @Param request body auth.authReq true "애플로그인 후 받은 id token"
// @Success 200 {object} auth.tokenRes
// @Router /auth/apple [post]
func (h *AppleHandler) VerifyAndIssueToken(c *gin.Context) {
	reqAuth := authReq{}
	if err := c.Bind(&reqAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding failure"})
		return
	}

	if reqAuth.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idToken is required"})
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
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("failed to get user by email: %v", err.Error())})
		return
	}

	if user == nil {
		c.JSON(http.StatusOK, tokenRes{
			AccessToken:  "",
			RefreshToken: "",
		})
		return
	}

	// 주섬토큰 발급
	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
