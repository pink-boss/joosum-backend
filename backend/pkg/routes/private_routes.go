package routes

import (
	"joosum-backend/app/tag"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {

	tagHandler := tag.TagHandler{}

	router.Use(middleware.SetUserData())

	tagRouter := router.Group("/tag")
	{
		tagRouter.GET("/", tagHandler.GetTags)
		tagRouter.POST("/", tagHandler.CreateTag)
		tagRouter.DELETE("/:id", tagHandler.DeleteTag)
	}

}
