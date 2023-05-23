package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"joosum-backend/docs"
	"strings"
)

// SwaggerRoutes func for describe group of Swagger routes.
func SwaggerRoutes(router *gin.Engine) {
	router.Use(getPrefix())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Host 에 따라 문서 basePath 설정
func getPrefix() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		basePath := "/api"

		if strings.HasPrefix(host, "localhost") {
			basePath = ""
		}
		docs.SwaggerInfo.BasePath = basePath
		c.Next()
	}
}
