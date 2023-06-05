package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/tag"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {

	tagHandler := tag.TagHandler{}
	authHandler := auth.AuthHandler{}

	router.Use(middleware.SetUserData())

	router.GET("/protected", authHandler.Protected)

	tagRouter := router.Group("/tags")
	{
		tagRouter.GET("/", tagHandler.GetTags)
		tagRouter.POST("/", tagHandler.CreateTag)
		tagRouter.DELETE("/:id", tagHandler.DeleteTag)
	}

}
