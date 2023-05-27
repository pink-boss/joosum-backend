package routes

import (
	"joosum-backend/app/tag"
	"joosum-backend/app/user"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {

	tagHandler := tag.TagHandler{}

	router.Use(middleware.SetUserData())

	// for apple test
	appleRouter := router.Group("/apple")
	{
		appleRouter.Use(middleware.AppleAuthMiddleware())
		router.GET("/protected", user.Protected)
	}
	

	tagRouter := router.Group("/tag")
	{
		tagRouter.GET("/", tagHandler.GetTags)
		tagRouter.POST("/", tagHandler.CreateTag)
		tagRouter.DELETE("/:id", tagHandler.DeleteTag)
	}

}
