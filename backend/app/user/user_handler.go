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
