package repository

import "auth-server/internal/domain"

type OAuthClientRepository interface {
	GetByID(id string) (*domain.OAuthClient, error)
}

type AuthCodeRepository interface {
	Save(code *domain.AuthCode) error
	Get(code string) (*domain.AuthCode, error)
	Delete(code string) error
}
