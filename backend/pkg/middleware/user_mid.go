package middleware

import (
	"joosum-backend/app/user"
	"joosum-backend/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var userUsecase user.UserUsecase  
func SetUserData() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header value
		tokenString := c.GetHeader("Authorization")

		// Initialize a new instance of `Claims`
		claims := jwt.MapClaims{}

		// Parse JWT string and store the result in `claims`.
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
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

		idValue, exists := claims["id"]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userId, err := userUsecase.GetUserByEmail(idValue.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", userId)

		c.Next()
		
	}
}