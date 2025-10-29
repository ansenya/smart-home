package services

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct {
}

func (s *passwordService) IsPasswordValid(password string) bool {
	return len(password) >= 8
}

func (s *passwordService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *passwordService) IsPasswordCorrect(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewPasswordService() PasswordService {
	return &passwordService{}
}
