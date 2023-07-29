package middleware

import (
	"fmt"
	"joosum-backend/app/user"
	"joosum-backend/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var userUsecase user.UserUsecase

func SetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header value
		tokenString := c.GetHeader("Authorization")
		TrimmedTkStr := strings.TrimPrefix(tokenString, "Bearer ")

		// Initialize a new instance of `Claims`
		claims := jwt.MapClaims{}

		// Parse JWT string and store the result in `claims`.
		token, err := jwt.ParseWithClaims(TrimmedTkStr, &claims, func(token *jwt.Token) (interface{}, error) {
			// Make sure that the token method conforms to `SigningMethodHMAC`
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.GetEnvConfig("jwt_secret")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		idValue, exists := claims["email"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		user, err := userUsecase.GetUserByEmail(idValue.(string))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "failed to find the email that Signed up"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("unauthorized : %v", err.Error())})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
