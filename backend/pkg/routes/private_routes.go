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
	linkHandler := link.LinkHandler{}

	router.Use(middleware.SetUserData())

	router.GET("/protected", authHandler.Protected)

	authRouter := router.Group("/auth")
	{
		authRouter.POST("/me", authHandler.GetMe)
	}
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
	}

	linkRouter := router.Group("/links")
	{
		linkRouter.POST("/", linkHandler.CreateLink)
		linkRouter.GET("/", linkHandler.GetLinks)
		linkRouter.GET("/:linkId", linkHandler.GetLinkByLinkId)
		linkRouter.GET("/link-books/:linkBookId/links", linkHandler.GetLinksByLinkBookId)
		linkRouter.DELETE("/:linkId", linkHandler.DeleteLinkByLinkId)
		linkRouter.DELETE("/", linkHandler.DeleteLinksByUserId)
		linkRouter.DELETE("/link-books/:linkBookId/links", linkHandler.DeleteLinksByLinkBookId)
		linkRouter.PUT("/:linkId/read-count", linkHandler.UpdateReadCount)
		linkRouter.PUT("/:linkId/link-book-id/:linkBookId", linkHandler.UpdateLinkBookIdByLinkId)
		linkRouter.PUT("/:linkId", linkHandler.UpdateTitleAndUrlByLinkId)
	}

}
