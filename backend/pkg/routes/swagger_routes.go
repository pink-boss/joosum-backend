package routes

import (
	"joosum-backend/docs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SwaggerRoutes func for describe group of Swagger routes.
func SwaggerRoutes(router *gin.Engine) {
	router.Use(getPrefix())

	// 스웨거에 배포를 한 날 (서버를 띄운 날) 이 노출
	dateServerStarted := time.Now().Format("060102")
	docs.SwaggerInfo.Version = dateServerStarted

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,

		// 스웨거 설정 - https://github.com/swaggo/gin-swagger
		ginSwagger.DocExpansion("nothing"), // Controls the default expansion setting for the operations and tags. It can be 'list' (expands only the tags), 'full' (expands the tags and operations) or 'none' (expands nothing).
		//ginSwagger.DeepLinking(),                // If set to true, enables deep linking for tags and operations. See the Deep Linking documentation for more information.
		//ginSwagger.DefaultModelsExpandDepth(), // 	Default expansion depth for models (set to -1 completely hide the models).
		//ginSwagger.InstanceName(), // The instance name of the swagger document. If multiple different swagger instances should be deployed on one gin router, ensure that each instance has a unique name (use the --instanceName parameter to generate swagger documents with swag init).
		//ginSwagger.PersistAuthorization(true), // 	If set to true, it persists authorization data and it would not be lost on browser close/refresh.
		//ginSwagger.Oauth2DefaultClientID(),      // If set, it's used to prepopulate the client_id field of the OAuth2 Authorization dialog.
	))
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
