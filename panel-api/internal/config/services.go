package config

import (
	"panel-api/internal/services"
)

type Services struct {
	OauthService services.OauthService
	UsersService services.UsersService
}
