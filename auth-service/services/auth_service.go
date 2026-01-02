package services

import (
	"auth-server/models"
	"auth-server/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	ErrorIncorrectPassword = errors.New("bad password")
)

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordService   PasswordService
	jwtService        JWTService
}

func (s *authService) Me(sid string) (*models.User, error) {
	session, err := s.sessionRepository.GetByID(sid)
	if err != nil || session == nil {
		return nil, fmt.Errorf("invalid session id")
	}
	user, err := s.userRepository.GetByID(session.UserID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("invalid session id")
	}
	return user, nil
}

func (s *authService) Login(request *models.AuthRequest) (*models.Session, error) {
	user, err := s.userRepository.GetByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if err := s.passwordService.IsPasswordCorrect(request.Password, user.Password); err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.jwtService.GetRefreshTokenDuration())
	session := models.Session{
		UserID:    user.ID,
		ExpiresAt: &expiresAt,
	}
	if err := s.sessionRepository.Create(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *authService) Register(request *models.AuthRequest) (*models.User, error) {
	if !s.passwordService.IsPasswordValid(request.Password) {
		return nil, ErrorIncorrectPassword
	}

	hash, err := s.passwordService.HashPassword(request.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, err
	}

	user := models.User{
		Email:    request.Email,
		Password: hash,
	}

	if err := s.userRepository.Create(&user); err != nil {
		log.Printf("error creating user: %v", err)
		return nil, err
	}
	return &user, nil
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		userRepository:    repository.NewUserRepository(db),
		sessionRepository: repository.NewSessionRepository(db),
		passwordService:   NewPasswordService(),
		jwtService:        NewJwtService(),
	}
}
