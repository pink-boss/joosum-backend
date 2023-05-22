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
	authUsecae   AuthUsecae
	googleUsecae GoogleUsecae
}

func (h GoogleHandler) VerifyGoogleAccessToken(c *gin.Context) {
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

	valid, err := h.googleUsecae.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if valid {
		email, err := h.googleUsecae.GetUserEmail(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := user.RegisterUser(email, "google")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := h.authUsecae.GenerateNewJWTToken([]string{"USER", "ADMIN"}, user.Email)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"email": user.Email,
			"accessToken":  accessToken,
			"refreshToken": refreshToken})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
	}
}
