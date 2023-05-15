package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// VerifyAppleAccessToken
// @Tags login
// @Summary 토큰 verify
// @Router /api/auth/apple [post]
func VerifyAppleAccessToken(c *gin.Context) {
	reqAuth := authRequest{}
	if err := c.Bind(&reqAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding failure"})
		return
	}

	if reqAuth.IdToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	claims, err := verifyToken(reqAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "success",
		"claims": claims,
	})
}

// GetAppleToken
// @Tags login
// @Summary access token, refresh 토큰 발급
// @Router /api/auth/apple/token [post]
func GetAppleToken(c *gin.Context) {
	reqAuth := authRequest{}
	res, err := getToken(reqAuth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"res": res,
	})
}
