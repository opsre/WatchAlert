package middleware

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"time"
)

// RequestLoggerFormatter 自定义的日志格式化函数
func RequestLoggerFormatter(param gin.LogFormatterParams) string {
	level := "info"
	switch {
	case param.StatusCode >= 500:
		level = "error"
	case param.StatusCode >= 400:
		level = "warn"
	case param.StatusCode >= 300:
		level = "debug"
	}

	logData := map[string]interface{}{
		"level":      level,
		"statusCode": param.StatusCode,
		"clientIP":   param.ClientIP,
		"method":     param.Method,
		"path":       param.Path,
		"time":       param.TimeStamp.Format(time.RFC3339),
	}

	// 将数据编码为JSON字符串
	jsonData, _ := sonic.Marshal(logData)
	return fmt.Sprintf("%s\n", jsonData)
}
