package routes

import (
	"joosum-backend/app/user"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// InternalRoutes 내부 서비스용 API 라우터를 설정합니다.
func InternalRoutes(router *gin.Engine) {
	userHandler := user.UserHandler{}

	internal := router.Group("")
	internal.Use(middleware.InternalAPIKeyMiddleware())
	{
		internal.GET("/withdraw-users", userHandler.GetWithdrawUsers)
		internal.GET("/signup-check", userHandler.CheckUserSignupByEmail)
	}
}
