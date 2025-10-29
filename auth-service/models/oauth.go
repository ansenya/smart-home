package models

type OauthRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	State    string `form:"state" binding:"required"`
	Scope    string `form:"scope"`
}
