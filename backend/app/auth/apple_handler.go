package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// IssueTokenFromApple
// @Tags 로그인
// @Summary id_token 을 verify 한 후 애플로 부터 토큰 발급
// @Param request body auth.authRequest true "code 와 id_token"
// @Success 200 {object} auth.tokenResponse
// @Router /api/auth/apple [post]
func IssueTokenFromApple(c *gin.Context) {
	reqAuth := authRequest{}
	if err := c.Bind(&reqAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding failure"})
		return
	}

	if reqAuth.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	resToken, err := issueTokenFromApple(reqAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  resToken.IdToken,
		RefreshToken: resToken.RefreshToken,
	})
}
