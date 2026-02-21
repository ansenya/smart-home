package repositories

import (
	"context"
	"llm-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	GetChatsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error)
	GetChatByID(ctx context.Context, id, userID uuid.UUID) (*models.Chat, error)
	UpdateChat(ctx context.Context, chat *models.Chat) error
	DeleteChat(ctx context.Context, id, userID uuid.UUID) error

	CreateMessage(ctx context.Context, msg *models.Message) error
	GetMessagesByChatID(ctx context.Context, chatID uuid.UUID, limit int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, msg *models.Message) error
}
type chatRepository struct {
	db *gorm.DB
}

func (r *chatRepository) CreateChat(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *chatRepository) GetChatsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error) {
	var chats []models.Chat

	err := r.db.WithContext(ctx).
		Where("user_id = ? and deleted_at IS NULL", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&chats).Error
	return chats, err
}

func (r *chatRepository) GetChatByID(ctx context.Context, id, userID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userID).
		First(&chat).Error
	return &chat, err
}

func (r *chatRepository) UpdateChat(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Save(chat).Error
}

func (r *chatRepository) DeleteChat(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Chat{}).Error
}

func (r *chatRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *chatRepository) GetMessagesByChatID(ctx context.Context, chatID uuid.UUID, limit int) ([]models.Message, error) {
	var msgs []models.Message
	err := r.db.WithContext(ctx).
		Where("chat_id = ? AND deleted_at IS NULL", chatID).
		Order("created_at ASC").
		Limit(limit).
		Find(&msgs).Error
	return msgs, err
}

func (r *chatRepository) UpdateMessage(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Save(msg).Error
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}
