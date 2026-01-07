package services

import (
	"log/slog"
	"panel-api/internal/models"
	"panel-api/internal/repositories"
	"time"
)

type UsersService interface {
	CreateSession(user *models.User, tokens *models.Tokens) (*models.Session, error)
	GetUserBySessionID(sid string) (*models.User, error)
	ExpireSession(sid string) error
}
type usersService struct {
	sessionRepo repositories.SessionRepository
	usersRepo   repositories.UsersRepository
	log         *slog.Logger
}

func NewUsersService(log *slog.Logger, sessionRepo repositories.SessionRepository, usersRepo repositories.UsersRepository) UsersService {
	return &usersService{
		sessionRepo: sessionRepo,
		usersRepo:   usersRepo,
		log:         log,
	}
}

func (s *usersService) CreateSession(user *models.User, tokens *models.Tokens) (*models.Session, error) {
	session := models.Session{
		UserID:       user.ID,
		TokenType:    tokens.TokenType,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return &session, s.sessionRepo.Create(&session)
}

func (s *usersService) GetUserBySessionID(sid string) (*models.User, error) {
	session, err := s.sessionRepo.Get(sid)
	if err != nil {
		return nil, err
	}

	return s.usersRepo.GetByID(session.UserID)
}

func (s *usersService) ExpireSession(sid string) error {
	session, err := s.sessionRepo.Get(sid)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now()
	return s.sessionRepo.Update(session)
}
