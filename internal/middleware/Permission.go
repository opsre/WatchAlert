package middleware

import (
	"fmt"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
	utils2 "watchAlert/pkg/tools"

	"github.com/bytedance/sonic"
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
		// 从上下文获取用户ID，支持JWT和API Key两种认证方式
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

		c := ctx.DO()
		// 获取当前用户
		var user models.Member
		err := c.DB.DB().Model(&models.Member{}).Where("user_id = ?", userId).First(&user).Error
		if gorm.ErrRecordNotFound == err {
			logc.Errorf(c.Ctx, fmt.Sprintf("用户不存在, uid: %s", userId))
		}
		if err != nil {
			response.PermissionFail(context)
			context.Abort()
			return
		}

		context.Set("UserEmail", user.Email)

		// 获取租户用户角色
		tenantUserInfo, tenantErr := c.DB.Tenant().GetTenantLinkedUserInfo(tid, userId)
		if tenantErr != nil {
			logc.Errorf(c.Ctx, fmt.Sprintf("获取租户用户角色失败 %s", tenantErr.Error()))
			response.TokenFail(context)
			context.Abort()
			return
		}
		if err != nil {
			response.PermissionFail(context)
			context.Abort()
			return
		}

		var (
			role       models.UserRole
			permission []models.UserPermissions
		)
		// 根据用户角色获取权限
		err = c.DB.DB().Model(&models.UserRole{}).Where("id = ?", tenantUserInfo.UserRole).First(&role).Error
		if err != nil {
			response.Fail(context, fmt.Sprintf("获取用户 %s 的角色失败, %s %s", user.UserName, tenantUserInfo.UserRole, err.Error()), "failed")
			logc.Errorf(c.Ctx, fmt.Sprintf("获取用户 %s 的角色失败 %s %s", user.UserName, tenantUserInfo.UserRole, err.Error()))
			context.Abort()
			return
		}
		_ = sonic.Unmarshal([]byte(utils2.JsonMarshalToString(role.Permissions)), &permission)

		urlPath := context.Request.URL.Path

		var pass bool
		for _, v := range permission {
			if urlPath == v.API {
				pass = true
				break
			}
		}
		if !pass {
			response.PermissionFail(context)
			context.Abort()
			return
		}
	}
}
