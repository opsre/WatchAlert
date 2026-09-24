package types

const OidcPassword = "watchAlert"

type RequestOidcCodeQuery struct {
	Code string `json:"code" form:"code"`
}

type OidcInfo struct {
	AuthType    *int   `json:"authType"`
	ClientID    string `json:"clientID"`
	UpperURI    string `json:"upperURI"`
	RedirectURI string `json:"redirectURI"`
}

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RespOpenIDConfiguration struct {
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	ClaimsSupported                   []string `json:"claims_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	IdTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	Issuer                            string   `json:"issuer"`
	JwksUri                           string   `json:"jwks_uri"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
}

type RespOidcUserInfo struct {
	Attributes struct {
		AvatarUrl   string        `json:"avatar_url"`
		Departments []interface{} `json:"departments"`
		Email       string        `json:"email"`
		Name        string        `json:"name"`
		Nickname    string        `json:"nickname"`
		PhoneNum    string        `json:"phone_num"`
	} `json:"attributes"`
	ClientId string `json:"client_id"`
	Email    string `json:"email"`
	Id       string `json:"id"`
	Sub      string `json:"sub"`
}
