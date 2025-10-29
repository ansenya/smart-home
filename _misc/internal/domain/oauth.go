package domain

type OAuthClient struct {
	ID          string
	Secret      string
	RedirectURI string
	Enabled     bool
}

type AuthCode struct {
	Code     string
	UserID   string
	ClientID string
	Scopes   []string
}
