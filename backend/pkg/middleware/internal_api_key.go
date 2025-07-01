package middleware

import (
	"net/http"

	"joosum-backend/pkg/config"

	"github.com/gin-gonic/gin"
)

// InternalAPIKeyMiddleware 는 내부 API 호출 시 사용되는 API 키를 검증한다.
// 환경변수 "internalApiKey" 가 설정되어 있어야 하며,
// 미설정 시 모든 요청을 거부한다.
func InternalAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Internal-Api-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			return
		}

		expected := config.GetEnvConfig("internalApiKey")
		if expected == "" {
			// 환경변수 설정이 없으면 기본 키를 사용하지 않고 요청을 거부한다
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		if key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}
