package services

import (
	"auth-server/models"
	"auth-server/repository"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type authService struct {
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
	passwordService   PasswordService
	jwtService        JWTService
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
		return nil, fmt.Errorf("bad password")
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
