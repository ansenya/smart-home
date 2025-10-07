package service

import (
	"auth-server/repository"
	"crypto/rsa"
	"time"
)

type AuthService struct {
	repo *repository.UserRepo

	smtpConfig SmtpConfig

	accessPrivateKey  *rsa.PrivateKey
	accessPublicKey   *rsa.PublicKey
	refreshPrivateKey *rsa.PrivateKey
	refreshPublicKey  *rsa.PublicKey

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type SmtpConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewAuthService(r *repository.UserRepo, smtpConfig SmtpConfig) (*AuthService, error) {
	accessPrivateKey, err := loadPrivateKey("keys/access_private.pem")
	if err != nil {
		return nil, err
	}

	accessPublicKey, err := loadPublicKey("keys/access_public.pem")
	if err != nil {
		return nil, err
	}

	refreshPrivateKey, err := loadPrivateKey("keys/refresh_private.pem")
	if err != nil {
		return nil, err
	}

	refreshPublicKey, err := loadPublicKey("keys/refresh_public.pem")
	if err != nil {
		return nil, err
	}

	return &AuthService{
		repo:                 r,
		smtpConfig:           smtpConfig,
		accessPrivateKey:     accessPrivateKey,
		accessPublicKey:      accessPublicKey,
		refreshPrivateKey:    refreshPrivateKey,
		refreshPublicKey:     refreshPublicKey,
		accessTokenDuration:  15 * time.Minute,
		refreshTokenDuration: 7 * 24 * time.Hour,
	}, nil
}
