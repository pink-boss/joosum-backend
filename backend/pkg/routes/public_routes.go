package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/user"

	"github.com/gin-gonic/gin"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(router *gin.Engine) {

	router.GET("/", user.GetMainPage)
	googleHandler := auth.GoogleHandler{}

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/google", googleHandler.VerifyGoogleAccessToken)
		authRouter.POST("/apple", auth.VerifyAppleAccessToken)
		authRouter.POST("/apple/token", auth.GetAppleToken)
	}

}
