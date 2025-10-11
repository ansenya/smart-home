package services

import (
	"auth-server/models"
	"auth-server/repository"
	"golang.org/x/crypto/bcrypt"
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

func (s *userService) IsPasswordValid(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func (s *userService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *userService) IsPasswordCorrect(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
