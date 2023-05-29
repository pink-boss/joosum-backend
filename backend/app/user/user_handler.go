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

// Protected
// @Tags 테스트
// @Summary 애플 JWT 미들웨어 테스트 api
// @Success 200 {string} Protected route
// @Router /protected [get]
func Protected(c *gin.Context) {
	// Your protected route logic goes here
	c.JSON(http.StatusOK, gin.H{"message": "Protected route"})
}
