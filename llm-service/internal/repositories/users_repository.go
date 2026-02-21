package repositories

import (
	"context"
	"llm-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type usersRepository struct {
	db *gorm.DB
}

func (r *usersRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return &usersRepository{
		db: db,
	}
}
