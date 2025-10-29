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

type OauthService interface {
	Authorize(queries models.OAuthRequest, sid string) (string, error)
	GetByID(id string) (*models.OauthClient, error)
	GetByName(name string) (*models.OauthClient, error)
}

type AuthService interface {
	Login(request *models.AuthRequest) (*models.Session, error)
	Register(request *models.AuthRequest) (*models.User, error)
}

type PasswordService interface {
	IsPasswordValid(password string) bool
	HashPassword(password string) (string, error)
	IsPasswordCorrect(password string, hash string) error
}

type TemporaryCodeService interface {
	Save(code string, data string, expiresIn time.Duration) error
	Get(code string) (string, error)
	Delete(code string) error
}

type JWTService interface {
	GenerateAccessToken(user *models.User) (string, error)
	ValidateAccessToken(token string) (*Claims, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateRefreshToken(token string) (*Claims, error)

	GetAccessTokenDuration() time.Duration
	GetRefreshTokenDuration() time.Duration

	GenerateJwks() Jwks
}
