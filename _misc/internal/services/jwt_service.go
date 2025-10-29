package services

import (
	"auth-server/internal/domain"
	"time"
)

type JWTService interface {
	GenerateAccessToken(*domain.User, []string) (string, error)
	GenerateRefreshToken(*domain.User) (string, error)
	Validate(token string) (*domain.User, error)
	AccessTokenTTL() time.Duration
}
