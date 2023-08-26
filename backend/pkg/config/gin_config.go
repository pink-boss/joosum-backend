package config

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/v3/log"
)

var Env *string

func GetRouter() *gin.Engine {
	Env = flag.String("env", "dev", "a string")
	flag.Parse()

	if *Env == "dev" {
		gin.SetMode(gin.DebugMode)
		return gin.Default()

	} else if *Env == "prod" {
		gin.SetMode(gin.ReleaseMode) // gin.New() 보다 위에 있어야 하는 듯
		router := gin.New()
		router.Use(gin.Recovery())
		return router

	} else {
		log.Errorf("Invalid env")
		return nil
	}
}
