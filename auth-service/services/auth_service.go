package services

import (
	"auth-server/models"
	"auth-server/repository"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordService   PasswordService
	jwtService        JWTService
}

func (s *authService) Login(request *models.LoginRequest) (*models.Session, error) {
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

func NewAuthService(db *gorm.DB) (AuthService, error) {
	jwtService, err := NewJwtService()
	if err != nil {
		return nil, fmt.Errorf("cannot initialize jwtService: %s", err)
	}
	return &authService{
		userRepository:    repository.NewUserRepository(db),
		sessionRepository: repository.NewSessionRepository(db),
		passwordService:   NewPasswordService(),
		jwtService:        jwtService,
	}, nil
}
