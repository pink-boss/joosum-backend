package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/link"
	"joosum-backend/app/tag"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {

	tagHandler := tag.TagHandler{}
	authHandler := auth.AuthHandler{}
	linkBookHandler := link.LinkBookHandler{}

	router.Use(middleware.SetUserData())

	router.GET("/protected", authHandler.Protected)

	tagRouter := router.Group("/tags")
	{
		tagRouter.GET("/", tagHandler.GetTags)
		tagRouter.POST("/", tagHandler.CreateTag)
		tagRouter.DELETE("/:id", tagHandler.DeleteTag)
	}

	linkBookRouter := router.Group("/link-books")
	{
		linkBookRouter.GET("/", linkBookHandler.GetLinkBooks)
		linkBookRouter.POST("/", linkBookHandler.CreateLinkBook)
		//linkBookRouter.DELETE("/:id", linkBookHandler.DeleteLinkBook)
	}

}
