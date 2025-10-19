package routes

import (
	"joosum-backend/app/auth"
	"joosum-backend/app/banner"
	"joosum-backend/app/link"
	"joosum-backend/app/notif"
	"joosum-backend/app/page"
	"joosum-backend/app/setting"
	"joosum-backend/app/tag"
	"joosum-backend/app/user"
	"joosum-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(router *gin.Engine) {

	UserHandler := user.UserHandler{}
	tagHandler := tag.TagHandler{}
	authHandler := auth.AuthHandler{}
	linkBookHandler := link.LinkBookHandler{}
	linkHandler := link.LinkHandler{}
	pageHandler := page.PageHandler{}
	settingHandler := setting.SettingHandler{}
	notificationHandler := notif.NotificationHandler{}
	BannerHandler := banner.BannerHandler{}

	router.Use(middleware.SetUserData())

	router.GET("/protected", authHandler.Protected)

	pageRouter := router.Group("/pages")
	{
		pageRouter.GET("/main", pageHandler.GetMainPage)
	}

	authRouter := router.Group("/auth")
	{
		authRouter.GET("/me", authHandler.GetMe)
		authRouter.DELETE("/me", UserHandler.DeleteUser)
		authRouter.POST("/logout", authHandler.Logout)
	}
	tagRouter := router.Group("/tags")
	{
		tagRouter.GET("", tagHandler.GetTags)
		tagRouter.POST("", tagHandler.CreateTags)
		tagRouter.DELETE("/:tag", tagHandler.DeleteTag)
	}

	linkBookRouter := router.Group("/link-books")
	{
		linkBookRouter.GET("", linkBookHandler.GetLinkBooks)
		linkBookRouter.POST("", linkBookHandler.CreateLinkBook)
		linkBookRouter.PUT("/:linkBookId", linkBookHandler.UpdateLinkBook)
		linkBookRouter.DELETE("/:linkBookId", linkBookHandler.DeleteLinkBook)
		linkBookRouter.GET("/:linkBookId/links", linkHandler.GetLinksByLinkBookId)
		linkBookRouter.DELETE("/:linkBookId/links", linkHandler.DeleteLinksByLinkBookId)
	}

	linkRouter := router.Group("/links")
	{
		linkRouter.POST("", linkHandler.CreateLink)
		linkRouter.GET("", linkHandler.GetLinks)
		linkRouter.GET("/:linkId", linkHandler.GetLinkByLinkId)
		linkRouter.DELETE("/:linkId", linkHandler.DeleteLinkByLinkId)
		linkRouter.DELETE("", linkHandler.DeleteLinksByUserId)
		linkRouter.PUT("/:linkId/read-count", linkHandler.UpdateReadCount)
		linkRouter.PUT("/:linkId/link-book-id/:linkBookId", linkHandler.UpdateLinkBookIdByLinkId)
		linkRouter.PUT("/:linkId", linkHandler.UpdateTitleAndUrlByLinkId)
		linkRouter.POST("/thumbnail", linkHandler.GetThumnailURL)
		linkRouter.POST("/ai-tags", linkHandler.GetAIRecommendedTags)
	}

	settingRouter := router.Group("/settings")
	{
		settingRouter.POST("/device", settingHandler.SaveDeviceId)
		settingRouter.GET("/notification", settingHandler.GetNotificationAgree)
		settingRouter.PUT("/notification", settingHandler.UpdatePushNotification)
	}

	notificationRouter := router.Group("/notifications")
	{
		notificationRouter.GET("", notificationHandler.Notifications)
		notificationRouter.PUT("/:notificationId", notificationHandler.ReadNotification)
	}

	bannerRouter := router.Group("/banners")
	{
		bannerRouter.GET("", BannerHandler.GetBanners)
		bannerRouter.POST("", BannerHandler.CreateBanner)
		bannerRouter.DELETE("/:bannerId", BannerHandler.DeleteBanner)
	}
}
