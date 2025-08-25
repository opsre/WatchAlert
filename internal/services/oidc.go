package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	if setting.OidcConfig.Enable != true {
		return &types.OidcInfo{
			Enable: false,
		}, nil
	}

	return &types.OidcInfo{
		Enable:      true,
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

	if setting.OidcConfig.Enable != true {
		return nil, fmt.Errorf("oidc is not enabled")
	}

	r := req.(*types.RequestOidcCodeQuery)

	data, err := oidc.GetOauthToken(setting.OidcConfig.UpperURI, r.Code)
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
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	result, err := oidc.GetCurrentUser(setting.OidcConfig.UpperURI, payload["uid"])
	if err != nil {
		return nil, err
	}

	_, ok, _ := os.ctx.DB.User().Get("", result.Data.BaseInfo.Name, "")
	if ok {
		logc.Infof(os.ctx.Ctx, fmt.Sprintf("用户 %s 已存在", result.Data.BaseInfo.Name))
	} else {
		err = os.ctx.DB.User().Create(models.Member{
			UserId:   tools.RandUid(),
			UserName: result.Data.BaseInfo.Name,
			Email:    result.Data.BaseInfo.Email,
			Phone:    result.Data.BaseInfo.PhoneNum,
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

	if setting.OidcConfig.Enable != true {
		return nil, fmt.Errorf("oidc is not enabled")
	}

	accessToken, err := ctx.Cookie("token")
	if err != nil {
		return nil, err
	}

	rfToken, err := ctx.Cookie("rftoken")
	if err != nil {
		return nil, err
	}

	if err = oidc.DecodeToken(setting.OidcConfig.UpperURI, accessToken, rfToken); err != nil {
		return nil, err
	}

	parts := strings.Split(accessToken, ".")
	if len(parts) < 2 {
		return nil, err
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err = json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	currentUser, err := oidc.GetCurrentUser(setting.OidcConfig.UpperURI, payload["uid"])
	if err != nil {
		return nil, err
	}

	data, _, err := os.ctx.DB.User().Get("", currentUser.Data.BaseInfo.Name, "")
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
	os.ctx.Redis.Redis().Set("uid-"+data.UserId, tools.JsonMarshal(r), duration)

	return models.ResponseLoginInfo{
		Token:    tokenData,
		Username: data.UserName,
		UserId:   data.UserId,
	}, nil
}
