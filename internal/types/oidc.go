package types

const OidcPassword = "watchalert"

type RequestOidcCodeQuery struct {
	Code string `json:"code" form:"code"`
}

type OidcInfo struct {
	Enable      bool   `json:"enable"`
	ClientID    string `json:"clientID"`
	UpperURI    string `json:"upperURI"`
	RedirectURI string `json:"redirectURI"`
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RespOidcAuth struct {
	Code   int  `json:"code"`
	Status bool `json:"status"`
	Data   struct {
		NewToken string `json:"new_token"`
	} `json:"data"`
}

type RespOidcUserInfo struct {
	Code int `json:"code"`
	Data struct {
		BaseInfo struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Nickname string `json:"nickname"`
			PhoneNum string `json:"phone_num"`
		} `json:"base_info"`
	} `json:"data"`
}
