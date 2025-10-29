package models

type OAuthRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	State    string `form:"state" binding:"required"`
	Scope    string `form:"scope"`
}
type OAuthData struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
}
