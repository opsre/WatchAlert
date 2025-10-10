package services

import (
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"strings"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/global"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/oidc"
	"watchAlert/pkg/tools"
)

type oidcService struct {
	ctx *ctx.Context
}

type InterOidcService interface {
	GetOidcInfo() (interface{}, interface{})
	CallBack(ctx *gin.Context, req interface{}) (interface{}, interface{})
	CookieConvertToken(ctx *gin.Context) (interface{}, interface{})
}

func newInterOidcService(ctx *ctx.Context) InterOidcService {
	return &oidcService{
		ctx: ctx,
	}
}

func (os oidcService) GetOidcInfo() (interface{}, interface{}) {
	setting, err := os.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}

	return &types.OidcInfo{
		AuthType:    setting.AuthType,
		ClientID:    setting.OidcConfig.ClientID,
		UpperURI:    setting.OidcConfig.UpperURI,
		RedirectURI: setting.OidcConfig.RedirectURI,
	}, nil
}

func (os oidcService) CallBack(ctx *gin.Context, req interface{}) (interface{}, interface{}) {
	setting, err := os.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}

	cfg, err := oidc.GetOpenIDConfiguration(setting.OidcConfig.UpperURI)
	if err != nil {
		return nil, err
	}

	r := req.(*types.RequestOidcCodeQuery)
	data, err := oidc.GetOauthToken(cfg.TokenEndpoint, r.Code)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(data.AccessToken, ".")
	if len(parts) < 2 {
		return nil, err
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err = sonic.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	result, err := oidc.GetCurrentUser(cfg.UserinfoEndpoint, data.AccessToken)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	_, ok, _ := os.ctx.DB.User().Get("", result.Id, "")
	if ok {
		logc.Infof(os.ctx.Ctx, fmt.Sprintf("用户 %s 已存在", result.Id))
	} else {
		err = os.ctx.DB.User().Create(models.Member{
			UserId:   tools.RandUid(),
			UserName: result.Id,
			Email:    result.Email,
			Phone:    result.Attributes.PhoneNum,
			Password: tools.GenerateHashPassword(types.OidcPassword),
			CreateBy: "OIDC",
			CreateAt: time.Now().Unix(),
		})
		if err != nil {
			return nil, err
		}
	}

	ctx.SetCookie("token", data.AccessToken, 60*60*24, "/", setting.OidcConfig.Domain, false, false)
	ctx.SetCookie("rftoken", data.RefreshToken, 60*60*24, "/", setting.OidcConfig.Domain, false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")

	return nil, nil
}

func (os oidcService) CookieConvertToken(ctx *gin.Context) (interface{}, interface{}) {
	setting, err := os.ctx.DB.Setting().Get()
	if err != nil {
		return nil, err
	}

	accessToken, err := ctx.Cookie("token")
	if err != nil {
		return nil, err
	}

	cfg, err := oidc.GetOpenIDConfiguration(setting.OidcConfig.UpperURI)
	if err != nil {
		return nil, err
	}

	result, err := oidc.GetCurrentUser(cfg.UserinfoEndpoint, accessToken)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	data, _, err := os.ctx.DB.User().Get("", result.Id, "")
	if err != nil {
		return nil, err
	}

	tokenData, err := tools.GenerateToken(data.UserId, data.UserName, data.Password)
	if err != nil {
		return nil, err
	}

	r := &types.RequestUserLogin{
		UserName: data.UserName,
		Email:    data.Email,
		Phone:    data.Phone,
		Password: tools.GenerateHashPassword(types.OidcPassword),
	}

	duration := time.Duration(global.Config.Jwt.Expire) * time.Second
	os.ctx.Redis.Redis().Set("uid-"+data.UserId, tools.JsonMarshalToString(r), duration)

	return models.ResponseLoginInfo{
		Token:    tokenData,
		Username: data.UserName,
		UserId:   data.UserId,
	}, nil
}
