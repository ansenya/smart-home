package models

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Email        string `json:"-"`
	Confirmed    bool   `json:"confirmed"`
	Name         string `json:"name"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginData struct {
	ClientID string `json:"clientId"`
	UserID   string `json:"userId"`
	State    string `json:"state"`
	Scope    string `json:"scope"`
}

type TokenRequest struct {
	Code         string `form:"code" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id" binding:"required"`
	RedirectURI  string `form:"redirect_uri" binding:"required"`
}
