package structs

type AuthData struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
}
