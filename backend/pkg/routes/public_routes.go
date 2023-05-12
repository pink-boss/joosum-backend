package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/common"
	"joosum-backend/app/user"

	"github.com/gin-gonic/gin"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(router *gin.Engine) {

	router.GET("/docs", common.GetDocs)
	router.GET("/", user.GetMainPage)

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/google", auth.VerifyGoogleAccessToken)
		authRouter.POST("/apple", auth.VerifyAppleAccessToken)
		authRouter.POST("/apple/token", auth.GetAppleToken)
	}

}
