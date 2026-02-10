package middleware

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

func Permission() gin.HandlerFunc {
	return func(context *gin.Context) {
		tid := context.Request.Header.Get(TenantIDHeaderKey)
		if tid == "null" || tid == "" {
			return
		}

		// 获取用户ID并校验
		userIdValue, exists := context.Get("UserId")
		if !exists {
			response.TokenFail(context)
			context.Abort()
			return
		}

		userId, ok := userIdValue.(string)
		if !ok {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 超级管理员免校验通道
		if userId == "admin" {
			context.Next()
			return
		}

		c := ctx.DO()

		// 获取当前用户信息
		var user models.Member
		err := c.DB.DB().Model(&models.Member{}).Where("user_id = ?", userId).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				logc.Errorf(c.Ctx, "用户不存在, uid: %s", userId)
			}
			response.PermissionFail(context)
			context.Abort()
			return
		}

		context.Set("UserEmail", user.Email)

		// 获取租户用户关联信息
		tenantUserInfo, err := c.DB.Tenant().GetTenantLinkedUserInfo(tid, userId)
		if err != nil {
			logc.Errorf(c.Ctx, "获取租户用户角色失败: %s", err.Error())
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 获取角色权限
		var role models.UserRole
		err = c.DB.DB().Model(&models.UserRole{}).Where("id = ?", tenantUserInfo.UserRole).First(&role).Error
		if err != nil {
			errMsg := fmt.Sprintf("获取用户 %s 的角色失败: %s", user.UserName, err.Error())
			logc.Errorf(c.Ctx, errMsg)
			response.Fail(context, errMsg, "failed")
			context.Abort()
			return
		}

		// 权限匹配
		urlPath := context.Request.URL.Path
		if !hasPermission(role.Permissions, urlPath) {
			response.PermissionFail(context)
			context.Abort()
			return
		}

		context.Next()
	}
}

// 权限匹配
func hasPermission(permissions []models.UserPermissions, currentPath string) bool {
	for _, v := range permissions {
		if currentPath == v.API {
			return true
		}
	}
	return false
}
