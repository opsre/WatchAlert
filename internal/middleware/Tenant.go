package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
)

const TenantIDHeaderKey = "TenantID"

func ParseTenant() gin.HandlerFunc {
	// 从HTTP头部获取TenantID并存储到上下文中，可以提高代码的可维护性、可重用性、安全性和性能，同时也使得错误处理和业务逻辑的实现更加高效和灵活。
	return func(context *gin.Context) {
		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "" {
			response.Fail(context, "租户ID不能为空", "failed")
			context.Abort()
			return
		}

		c := ctx.DO()

		var count int64
		err := c.DB.DB().Model(&models.Tenant{}).Where("id = ?", tid).Count(&count).Error

		if count == 0 {
			response.Fail(context, "租户不存在", "failed")
			context.Abort()
			return
		}

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Fail(context, "租户不存在", "failed")
			} else {
				response.Fail(context, "数据库查询失败: "+err.Error(), "failed")
			}
			context.Abort()
			return
		}

		context.Set(TenantIDHeaderKey, tid)
		context.Next()
	}
}
