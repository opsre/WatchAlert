package tools

import (
	"errors"
	"time"
	"watchAlert/config"

	"github.com/dgrijalva/jwt-go"
)

// JwtCustomClaims 注册声明是JWT声明集的结构化版本，仅限于注册声明名称
type JwtCustomClaims struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Pass           string `json:"pass"`
	StandardClaims jwt.StandardClaims
}

const (
	// TokenType Token 类型
	TokenType = "bearer"
	// AppGuardName 颁发者
	AppGuardName = "WatchAlert"
)

// signKey 返回 JWT 签名密钥, 取自已加载的配置(config.Application.Jwt.Secret)。
func signKey() []byte {
	return []byte(config.Application.Jwt.Secret)
}

func (j JwtCustomClaims) Valid() error {
	if j.StandardClaims.ExpiresAt > 0 && time.Now().Unix() > j.StandardClaims.ExpiresAt {
		return errors.New("token is expired")
	}
	return nil
}

// ParseToken 解析token
func ParseToken(tokenStr string) (JwtCustomClaims, error) {
	iJwtCustomClaims := JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return signKey(), nil
	})

	if err == nil && !token.Valid {
		err = errors.New("invalid Token")
	}
	return iJwtCustomClaims, err
}

// GenerateToken 生成Token
func GenerateToken(userId, userName, password string) (string, error) {
	// 初始化
	iJwtCustomClaims := JwtCustomClaims{
		ID:   userId,
		Name: userName,
		Pass: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + config.Application.Jwt.Expire,
			IssuedAt:  time.Now().Unix(),
			Issuer:    AppGuardName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)
	return token.SignedString(signKey())
}

func GetUser(tokenStr string) string {
	if tokenStr == "" {
		return ""
	}

	tokenStr = tokenStr[len(TokenType)+1:]
	token, err := ParseToken(tokenStr)
	if err != nil {
		return ""
	}
	return token.Name
}

func GetUserID(tokenStr string) string {
	if tokenStr == "" {
		return ""
	}

	tokenStr = tokenStr[len(TokenType)+1:]
	token, err := ParseToken(tokenStr)
	if err != nil {
		return ""
	}

	return token.ID
}
