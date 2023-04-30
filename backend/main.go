package main

import (
	_ "github.com/create-go-app/net_http-go-template/docs" // load Swagger docs
	"github.com/create-go-app/net_http-go-template/pkg/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs for Golang net/http Template.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	//Initialize a new router.
	//router := mux.NewRouter()
	//
	//// List of app routes:
	//routes.PublicRoutes(router)
	//routes.PrivateRoutes(router)
	//routes.SwaggerRoutes(router)
	//
	//// Register middleware.
	//router.Use(mux.CORSMethodMiddleware(router)) // enable CORS
	//
	//// Initialize server.
	//server := configs.ServerConfig(router)
	//
	//// Start API server.
	//utils.StartServerWithGracefulShutdown(server)

	router := gin.Default()
	handler, _ := routes.NewHandler()

	router.GET("/", handler.GetMainPage)
	router.GET("/api/v1/users", handler.GetUsers)
	router.GET("/api/v1/users/:id", handler.GetUser)

	router.Run(":5001") // listen and serve on 0.0.0.0:5000 (for windows "localhost:5000")
}
