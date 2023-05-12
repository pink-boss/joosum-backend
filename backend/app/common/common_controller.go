package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDocs(c *gin.Context) {
	host := c.Request.Host

	if c.Request.TLS != nil {
		c.Redirect(http.StatusMovedPermanently, "https://"+host+"/swagger/index.html")
	} else {
		c.Redirect(http.StatusMovedPermanently, "http://"+host+"/swagger/index.html")
	}
}
