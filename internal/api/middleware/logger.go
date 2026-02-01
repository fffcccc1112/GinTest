package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test/pkg/logger"
	"time"
)

// 日志拦截器
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		cost := time.Since(startTime)
		// 记录日志
		logger.Info("请求日志",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("cost", cost),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
