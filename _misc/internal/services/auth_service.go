package services

import (
	"auth-server/internal/domain"
	"auth-server/internal/repository"
	"context"
	"errors"
)

type AuthService struct {
	users     repository.UserRepository
	oauthRepo repository.OAuthClientRepository
	codes     repository.AuthCodeRepository
	jwt       JWTService
	password  PasswordService
}

func NewAuthService(
	users repository.UserRepository,
	oauthRepo repository.OAuthClientRepository,
	codes repository.AuthCodeRepository,
	jwt JWTService,
	password PasswordService,
) *AuthService {
	return &AuthService{users, oauthRepo, codes, jwt, password}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if !s.password.Validate(password) {
		return nil, errors.New("weak_password")
	}
	hash, _ := s.password.Hash(password)
	user := &domain.User{Email: email, Password: hash}
	return user, s.users.Create(user)
}
