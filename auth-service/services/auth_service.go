package services

import (
	"auth-server/models"
	"auth-server/repository"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrorIncorrectPassword = errors.New("bad password")
	ErrorEmailExists       = errors.New("email already registered")
	ErrorInvalidEmail      = errors.New("invalid email")
	ErrorInvalidCredentials = errors.New("invalid email or password")
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
	if session.ExpiresAt != nil && session.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("session expired")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorInvalidCredentials
		}
		return nil, err
	}
	if user == nil {
		return nil, ErrorInvalidCredentials
	}

	if err := s.passwordService.IsPasswordCorrect(request.Password, user.Password); err != nil {
		return nil, ErrorInvalidCredentials
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

func (s *authService) Register(request *models.AuthRequest) (*models.Session, error) {
	email := strings.TrimSpace(strings.ToLower(request.Email))
	if email == "" || !strings.Contains(email, "@") {
		return nil, ErrorInvalidEmail
	}
	if !s.passwordService.IsPasswordValid(request.Password) {
		return nil, ErrorIncorrectPassword
	}

	hash, err := s.passwordService.HashPassword(request.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, err
	}

	user := models.User{
		Email:    email,
		Password: hash,
	}

	if err := s.userRepository.Create(&user); err != nil {
		if isDuplicateEmail(err) {
			return nil, ErrorEmailExists
		}
		log.Printf("error creating user: %v", err)
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

func (s *authService) Logout(sid string) error {
	if sid == "" {
		return nil
	}
	return s.sessionRepository.Delete(sid)
}

func isDuplicateEmail(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") && strings.Contains(msg, "users_email_key")
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		userRepository:    repository.NewUserRepository(db),
		sessionRepository: repository.NewSessionRepository(db),
		passwordService:   NewPasswordService(),
		jwtService:        NewJwtService(),
	}
}
