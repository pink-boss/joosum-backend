package middleware

import (
	"net/http"

	"joosum-backend/pkg/config"

	"github.com/gin-gonic/gin"
)

// InternalAPIKeyMiddleware 는 내부 API 호출 시 사용되는 API 키를 검증한다.
func InternalAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Api-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			return
		}

		expected := config.GetEnvConfig("internalApiKey")
		if expected == "" {
			expected = "test-internal-key"
		}

		if key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}
