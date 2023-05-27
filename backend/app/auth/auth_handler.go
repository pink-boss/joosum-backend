package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecae *AuthUsecae
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	email, social := req.Email, req.Social
	if email == "" || social == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and social are required"})
		return
	}

	user, err := h.authUsecae.SignUp(email, social)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
