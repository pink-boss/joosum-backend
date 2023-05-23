package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMainPage
// @Tags Main
// @Summary
// @Router / [get]
func GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
}

func GetUsers(c *gin.Context) {
	c.String(http.StatusOK, "GetUsers API!!")
}

func GetUser(c *gin.Context) {
	email := c.Param("email")
	user, err := GetUserByEmail(email)
	if err != nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	c.String(http.StatusOK, "CreateUser API!!")
}

func UpdateUser(c *gin.Context) {
	c.String(http.StatusOK, "UpdateUser API!!")
}

func DeleteUser(c *gin.Context) {
	c.String(http.StatusOK, "DeleteUser API!!")
}

// Protected
// @Tags 테스트
// @Summary 애플 JWT 미들웨어 테스트 api
// @Success 200 {string} Protected route
// @Router /protected [get]
func Protected(c *gin.Context) {
	// Your protected route logic goes here
	c.JSON(http.StatusOK, gin.H{"message": "Protected route"})
}
