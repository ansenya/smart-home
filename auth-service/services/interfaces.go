package services

import (
	"auth-server/models"
	"time"
)

type UserService interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	IsPasswordValid(password string) bool
	HashPassword(password string) (string, error)
	IsPasswordCorrect(password string, hash string) error
}

type OauthClientsService interface {
	GetByID(id string) (*models.OauthClient, error)
	GetByName(name string) (*models.OauthClient, error)
}

type TemporaryCodeService interface {
	Save(code string, data string, expiresIn time.Duration) error
	Get(code string) (string, error)
	Delete(code string) error
}

type JWTService interface {
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateToken(token string) (*Claims, error)

	GetAccessTokenDuration() time.Duration
	GetRefreshTokenDuration() time.Duration
}
