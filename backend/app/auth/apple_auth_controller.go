package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type authResponse struct {
	State   string
	Code    string
	IdToken string `json:"id_token"`
}

func VerifyAppleAccessToken(c *gin.Context) {
	reqAuth := authResponse{}
	if err := c.Bind(&reqAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding failure"})
		return
	}

	if reqAuth.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token is required"})
		return
	}

	res, err := getApplePublicKeys(reqAuth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	println(res)

	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
