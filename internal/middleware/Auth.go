package middleware

import (
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
	"watchAlert/pkg/tools"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
)

const ApiKeyHeader = "X-API-Key"

func Auth() gin.HandlerFunc {

	return func(context *gin.Context) {
		// 获取 Token
		tokenStr := context.Request.Header.Get("Authorization")
		apiKey := context.Request.Header.Get(ApiKeyHeader)

		// 优先检查 JWT Token
		if tokenStr != "" {
			// 校验 Token
			ok := IsTokenValid(ctx.DO(), tokenStr)
			if !ok {
				response.TokenFail(context)
				context.Abort()
				return
			}
			// JWT验证成功后，也需要从Token中提取用户ID并存储到上下文中
			userId := tools.GetUserID(tokenStr)
			context.Set("UserId", userId)
		} else if apiKey != "" {
			// 如果没有Token，则尝试API Key认证
			userId, ok := IsApiKeyValid(ctx.DO(), apiKey)
			if !ok {
				response.TokenFail(context)
				context.Abort()
				return
			}
			// 将用户ID存储到上下文中，供后续处理使用
			context.Set("UserId", userId)
		} else {
			// 如果两者都没有提供
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 继续执行后续处理器
		context.Next()
	}
}

func IsTokenValid(ctx *ctx.Context, tokenStr string) bool {
	// Bearer Token, 获取 Token 值
	tokenStr = tokenStr[len(tools.TokenType)+1:]
	token, err := tools.ParseToken(tokenStr)
	if err != nil {
		logc.Errorf(ctx.Ctx, "parse token error: %v", err)
		return false
	}

	// 发布者校验
	if token.StandardClaims.Issuer != tools.AppGuardName {
		return false
	}

	// 密码校验, 当修改密码后其他已登陆的终端会被下线。
	var user models.Member
	result, err := ctx.Redis.Redis().Get("uid-" + token.ID).Result()
	if err != nil {
		logc.Errorf(ctx.Ctx, "get user by id error: %v", err)
		return false
	}
	_ = sonic.Unmarshal([]byte(result), &user)

	if token.Pass != user.Password {
		return false
	}

	// 校验过期时间
	ok := token.StandardClaims.VerifyExpiresAt(time.Now().Unix(), false)
	if !ok {
		return false
	}

	return true

}

// IsApiKeyValid 验证API密钥的有效性
func IsApiKeyValid(ctx *ctx.Context, apiKey string) (string, bool) {
	// 查询数据库中是否存在该API密钥
	apiKeyModel, exists, err := ctx.DB.ApiKey().GetByKey(apiKey)
	if err != nil || !exists {
		return "无效的API密钥", false
	}

	// 返回用户ID和成功标志
	return apiKeyModel.UserId, true
}
