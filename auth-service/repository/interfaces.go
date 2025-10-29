package repository

import "auth-server/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

type OauthClientsRepository interface {
	GetByID(id string) (*models.OauthClient, error)
}

type SessionRepository interface {
	Create(session *models.Session) error
	GetByID(id string) (*models.Session, error)
}
