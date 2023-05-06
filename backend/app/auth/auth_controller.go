package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyGoogleAccessToken(c *gin.Context) {
	accessToken := c.PostForm("accessToken")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_token is required"})
		return
	}

	valid, err := ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if valid {
		c.JSON(http.StatusOK, gin.H{"valid": true})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
	}
}
