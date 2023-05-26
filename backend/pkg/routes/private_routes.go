package routes

import (
	"github.com/gin-gonic/gin"
	"joosum-backend/app/user"
	"joosum-backend/pkg/middleware"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {
	router.Use(middleware.AppleAuthMiddleware())

	// Protected route
	router.GET("/protected", user.Protected)
}
