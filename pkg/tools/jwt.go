package tools

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
	"watchAlert/internal/global"
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

func (j JwtCustomClaims) Valid() error {
	return nil
}

// ParseToken 解析token
func ParseToken(tokenStr string) (JwtCustomClaims, error) {
	iJwtCustomClaims := JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return global.StSignKey, nil
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
			ExpiresAt: time.Now().Unix() + global.Config.Jwt.Expire,
			IssuedAt:  time.Now().Unix(),
			Issuer:    AppGuardName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)
	return token.SignedString(global.StSignKey)
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
