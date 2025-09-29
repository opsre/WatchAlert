package middleware

import (
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/response"
	"watchAlert/pkg/tools"
)

func Auth() gin.HandlerFunc {

	return func(context *gin.Context) {
		// 获取 Token
		tokenStr := context.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(context)
			context.Abort()
			return
		}

		// 校验 Token
		code, ok := IsTokenValid(ctx.DO(), tokenStr)
		if !ok {
			if code == 401 {
				response.TokenFail(context)
				context.Abort()
				return
			}
		}

	}
}

func IsTokenValid(ctx *ctx.Context, tokenStr string) (int64, bool) {
	// Bearer Token, 获取 Token 值
	tokenStr = tokenStr[len(tools.TokenType)+1:]
	token, err := tools.ParseToken(tokenStr)
	if err != nil {
		return 400, false
	}

	// 发布者校验
	if token.StandardClaims.Issuer != tools.AppGuardName {
		return 400, false
	}

	// 密码校验, 当修改密码后其他已登陆的终端会被下线。
	var user models.Member
	result, err := ctx.Redis.Redis().Get("uid-" + token.ID).Result()
	if err != nil {
		return 400, false
	}
	_ = sonic.Unmarshal([]byte(result), &user)

	if token.Pass != user.Password {
		return 401, false
	}

	// 校验过期时间
	ok := token.StandardClaims.VerifyExpiresAt(time.Now().Unix(), false)
	if !ok {
		return 401, false
	}

	return 200, true

}
