package auth

import (
	"net/http"

	"joosum-backend/app/user"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *AuthUsecase
	userUsecase *user.UserUsecase
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	isExist, err := h.userUsecase.GetUserByEmail(req.Email); 
	if isExist != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	var userInfo user.User
	userInfo.Email = req.Email
	userInfo.Social = req.Social
	userInfo.Name = req.Nickname
	userInfo.Age = req.Age
	userInfo.Gender = req.Gender

	_, err = h.authUsecase.SignUp(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authUsecase.GenerateNewJWTToken([]string{"USER", "ADMIN"}, req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken , "refreshToken": refreshToken})
}
