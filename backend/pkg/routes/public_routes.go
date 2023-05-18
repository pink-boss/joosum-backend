package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/user"

	"github.com/gin-gonic/gin"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(router *gin.Engine) {

	router.GET("/", user.GetMainPage)

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/google", auth.VerifyGoogleAccessToken)
		authRouter.POST("/apple", auth.IssueTokenFromApple)
		authRouter.POST("/apple/callback", auth.GetTokenFromApple)
	}

}
