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
