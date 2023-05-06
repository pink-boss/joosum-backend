package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMainPage(c *gin.Context) {
	log.Println("Main page....")
	c.String(http.StatusOK, "Main page for secure API!!")
}

// GetUsers func gets all exists users.
// @Description Get all exists users.
// @Summary get all exists users
// @Tags Users
// @Accept json
// @Produce json
// @Router /v1/users [get]
func GetUsers(c *gin.Context) {
	c.String(http.StatusOK, "GetUsers API!!")
}

// GetUser func gets one user by given ID or 404 error.
// @Description Get user by given ID.
// @Summary get user by given ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Router /v1/user/{id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"id":    id,
		"name":  "홍길동",
		"email": "gildong@gmail.com",
	})
}

// CreateUser func for creates a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "E-mail"
// @Security ApiKeyAuth
// @Router /v1/user [post]
func CreateUser(c *gin.Context) {
	c.String(http.StatusOK, "CreateUser API!!")
}

// UpdateUser func for updates user by given ID.
// @Description Update user.
// @Summary update user
// @Tags User
// @Accept json
// @Produce json
// @Param id body string true "User ID"
// @Param user_status body integer true "User status"
// @Security ApiKeyAuth
// @Router /v1/user [put]
func UpdateUser(c *gin.Context) {
	c.String(http.StatusOK, "UpdateUser API!!")
}

// DeleteUser func for deletes user by given ID.
// @Description Delete user by given ID.
// @Summary delete user by given ID
// @Tags User
// @Accept json
// @Produce json
// @Param id body string true "User ID"
// @Success 200 {string} string "ok"
// @Security ApiKeyAuth
// @Router /v1/user [delete]
func DeleteUser(c *gin.Context) {
	c.String(http.StatusOK, "DeleteUser API!!")
}
