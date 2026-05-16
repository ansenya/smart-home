package repositories

import (
	"context"
	"llm-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Chat, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error)
	Update(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type chatRepository struct {
	db *gorm.DB
}

func (r *chatRepository) Create(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *chatRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&chat).Error
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error) {
	var chats []models.Chat
	q := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepository) Update(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

func (r *chatRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Chat{}).Error
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}
