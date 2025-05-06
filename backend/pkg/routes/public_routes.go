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
	authHandler := auth.AuthHandler{}
	appleHandler := auth.AppleHandler{}
	naverHandler := auth.NaverHandler{}

	authRouter := router.Group("/auth")
	{

		authRouter.GET("/google/web", googleHandler.AuthGoogleWebLogin)
		authRouter.GET("/google/callback", googleHandler.AuthGoogleWebCallback)
		authRouter.POST("/google/web", googleHandler.VerifyGoogleAccessTokenInWeb)
		authRouter.POST("/apple", appleHandler.VerifyAndIssueToken)
		authRouter.POST("/google", googleHandler.VerifyGoogleAccessToken)
		authRouter.POST("/google/android", googleHandler.VerifyGoogleAccessTokenInAndroid)
		authRouter.POST("/naver", naverHandler.VerifyNaverAccessToken)
		authRouter.POST("/signup", authHandler.SignUp)
	}

}
