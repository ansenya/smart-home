package services

import (
	"auth-server/models"
	"auth-server/repository"
	"gorm.io/gorm"
)

type userService struct {
	userRepository repository.UserRepository
}

func (s *userService) Create(user *models.User) error {
	return s.userRepository.Create(user)
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.userRepository.GetByID(id)
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.userRepository.GetByEmail(email)
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		userRepository: repository.NewUserRepository(db),
	}
}
