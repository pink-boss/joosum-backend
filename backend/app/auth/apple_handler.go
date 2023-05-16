package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// IssueTokenFromApple
// @Tags login
// @Summary id_token 을 verify 한 후 애플로 부터 토큰 발급
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

	token, err := issueTokenFromApple(reqAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    token,
	})
}
