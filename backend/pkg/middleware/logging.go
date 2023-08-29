package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"time"
)

var logger *zap.Logger

func init() {
	config := zap.NewProductionConfig()

	encoderConfig := zapcore.EncoderConfig{
		TimeKey: "timestamp", // ts -> timestamp
		//LevelKey:       "level", // info 미노출
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 1692812828.761468 -> 2023-08-24T02:47:51.690+0900
		EncodeDuration: zapcore.SecondsDurationEncoder, // elapsed_time (초 단위)
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config.EncoderConfig = encoderConfig
	logger, _ = config.Build() // 아무것도 넣지 않으면 caller 에 로깅위치 나옴
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Log the request
		requestBody, _ := ioutil.ReadAll(c.Request.Body)
		var jsonReq interface{}
		json.Unmarshal(requestBody, &jsonReq)

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("query", c.Request.URL.Query()),
			zap.Any("request_body", jsonReq),
		)

		// Start the timer
		start := time.Now()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		var jsonRes interface{}
		json.Unmarshal(blw.body.Bytes(), &jsonRes)

		// Log the response
		elapsed := time.Since(start)
		logger.Info("Response",
			zap.Int("status_code", c.Writer.Status()),
			zap.Any("response_body", jsonRes),
			zap.Duration("elapsed_time", elapsed),
		)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
