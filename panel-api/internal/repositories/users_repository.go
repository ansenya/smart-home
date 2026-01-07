package repositories

import (
	"panel-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	GetByID(id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	return &user, r.db.Where("id = ?", id).First(&user).Error
}
