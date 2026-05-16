package repositories

import (
	"context"
	"llm-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *models.Message) error
	ListByChat(ctx context.Context, chatID uuid.UUID, limit int) ([]models.Message, error)
	DeleteByChat(ctx context.Context, chatID uuid.UUID) error
	DeleteByChatExceptLast(ctx context.Context, chatID uuid.UUID, keep int) error
}

type messageRepository struct {
	db *gorm.DB
}

func (r *messageRepository) Create(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *messageRepository) ListByChat(ctx context.Context, chatID uuid.UUID, limit int) ([]models.Message, error) {
	var msgs []models.Message
	q := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *messageRepository) DeleteByChat(ctx context.Context, chatID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Delete(&models.Message{}).Error
}

func (r *messageRepository) DeleteByChatExceptLast(ctx context.Context, chatID uuid.UUID, keep int) error {
	if keep <= 0 {
		return r.DeleteByChat(ctx, chatID)
	}

	var ids []uuid.UUID
	if err := r.db.WithContext(ctx).
		Model(&models.Message{}).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(keep).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("chat_id = ? AND id NOT IN ?", chatID, ids).
		Delete(&models.Message{}).Error
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}
