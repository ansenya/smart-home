package models

type OauthRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	State    string `form:"state" binding:"required"`
	Scope    string `form:"scope"`
}

type RefreshTokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	RefreshToken string `form:"refresh_token" binding:"required"`
}

type AccessTokenRequest struct {
	Code         string `form:"code" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
}

type TokenResponse struct {
	TokenType    string  `json:"token_type"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    float64 `json:"expires_in"`
}

type OauthData struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
}
