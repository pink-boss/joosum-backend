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
