package services

import (
	"auth-server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type OauthService interface {
	Authorize(queries models.OauthRequest, sid string) (string, error)
	GenerateTokens(request models.AccessTokenRequest) (*models.TokenResponse, error)
	RefreshTokens(request models.RefreshTokenRequest) (*models.TokenResponse, error)
	GetUserinfo(accessToken string) (*models.User, error)
}

type AuthService interface {
	Me(sid string) (*models.User, error)
	Login(request *models.AuthRequest) (*models.Session, error)
	Register(request *models.AuthRequest) (*models.Session, error)
}

type JWTService interface {
	GenerateAccessToken(user *models.User) (string, error)
	ValidateAccessToken(token string) (*jwt.RegisteredClaims, error)
	GenerateRefreshToken(user *models.User) (string, error)
	ValidateRefreshToken(token string) (*jwt.RegisteredClaims, error)

	GetAccessTokenDuration() time.Duration
	GetRefreshTokenDuration() time.Duration

	GenerateJwks() Jwks
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
