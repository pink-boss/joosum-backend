package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

// LoggingMiddleware는 민감한 정보 노출을 막기 위해 요청과 응답 본문을 기록하지 않는다.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 요청 시작 시간 기록
		start := time.Now()

		c.Next()

		// 민감 정보 노출을 피하기 위해 본문은 기록하지 않는다
		elapsed := time.Since(start)
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Any("query", c.Request.URL.Query()),
			zap.Int("status_code", c.Writer.Status()),
			zap.Duration("elapsed_time", elapsed),
		)
	}
}
