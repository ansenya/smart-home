package services

import (
	"log/slog"
	"panel-api/internal/models"
	"panel-api/internal/repositories"
)

type UsersService interface {
	CreateSession(user *models.User, tokens *models.Tokens) (*models.Session, error)
}
type usersService struct {
	repo *repositories.SessionRepository
	log  *slog.Logger
}

func NewUsersService(log *slog.Logger, repo *repositories.SessionRepository) UsersService {
	return &usersService{
		repo: repo,
		log:  log,
	}
}

func (s *usersService) CreateSession(user *models.User, tokens *models.Tokens) (*models.Session, error) {
	session := models.Session{
		UserID:       user.ID,
		TokenType:    tokens.TokenType,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return &session, s.repo.Create(&session)
}
