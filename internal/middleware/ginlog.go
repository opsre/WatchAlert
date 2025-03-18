package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"time"
)

// GinZapLogger returns a gin.HandlerFunc that logs requests using zap
func GinZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()

		// 计算请求耗时
		latency := end.Sub(start)

		// 获取请求的相关信息
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		message := c.Errors.ByType(gin.ErrorTypePrivate).String()

		ctx := logx.ContextWithFields(context.Background(),
			logx.Field("method", method),
			logx.Field("path", path),
			logx.Field("status", status),
			logx.Field("clientIP", clientIP),
			logx.Field("latency", latency),
		)
		logc.Info(ctx, message)
	}
}

// LoggingMiddleware 打印请求的body和query params
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取body
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			logc.Errorf(context.Background(), err.Error())
			c.Abort()
			return
		}
		// 将body复制回原位
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// 打印body和query params
		fmt.Println("Body:", string(bodyBytes))
		fmt.Println("Query Params:", c.Request.URL.Query())
		// 处理请求
		c.Next()
	}
}
