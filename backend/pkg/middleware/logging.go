package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

var logger *zap.Logger

func init() {
	// Initialize the logger
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Log the request
		requestBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.ByteString("request_body", requestBody),
		)

		// Start the timer
		start := time.Now()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		// Log the response
		elapsed := time.Since(start)
		logger.Info("Response",
			zap.Int("status_code", c.Writer.Status()),
			zap.String("response_body", blw.body.String()),
			zap.Duration("elapsed_time", elapsed),
		)
	}
}
