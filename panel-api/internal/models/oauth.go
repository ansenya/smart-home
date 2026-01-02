package models

type CodeExchange struct {
	Code     string `json:"code"`
	ClientID string `form:"client_id" binding:"required"`
}

type AccessTokenRequest struct {
	Code         string `form:"code" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
}

type Tokens struct {
	TokenType    string  `json:"token_type"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    float64 `json:"expires_in"`
}
