package repository

import (
	"auth-server/models"
	"errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(&user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID, passwordHash string) error {
	res := r.db.Model(&models.User{}).Where("id = ?", userID).Update("password", passwordHash)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
