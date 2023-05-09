package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccessTokenRequest struct {
	AccessToken string `json:"accessToken"`
}

func VerifyGoogleAccessToken(c *gin.Context) {
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
