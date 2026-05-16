package session

import (
	"context"
	"llm-service/internal/clients"
	"llm-service/internal/models"
	"llm-service/internal/repositories"

	"github.com/google/uuid"
)

const (
	MessagesLimit = 50
)

type Store interface {
	Get(ctx context.Context, chatID uuid.UUID) ([]clients.Message, error)
	Append(ctx context.Context, chatID uuid.UUID, msg clients.Message) error
	Clear(ctx context.Context, chatID uuid.UUID) error
	Trim(ctx context.Context, chatID uuid.UUID, maxMessages int) error
}

type store struct {
	ChatRepository    repositories.ChatRepository
	MessageRepository repositories.MessageRepository
}

func (s *store) Get(ctx context.Context, chatID uuid.UUID) ([]clients.Message, error) {
	msgs, err := s.MessageRepository.ListByChat(ctx, chatID, MessagesLimit)
	if err != nil {
		return nil, err
	}

	out := make([]clients.Message, 0, len(msgs))
	for _, m := range msgs {
		cm := clients.Message{
			Role:    clients.Role(m.Role),
			Content: m.Content,
		}
		if m.ToolCallID != nil {
			cm.ToolCallID = *m.ToolCallID
		}
		if m.ToolName != nil {
			cm.Name = *m.ToolName
		}
		out = append(out, cm)
	}
	return out, nil
}

func (s *store) Append(ctx context.Context, chatID uuid.UUID, msg clients.Message) error {
	m := &models.Message{
		ChatID:  chatID,
		Role:    models.MessageRole(msg.Role),
		Content: msg.Content,
		Status:  models.StatusCompleted,
	}

	if msg.ToolCallID != "" {
		m.ToolCallID = &msg.ToolCallID
	}
	if msg.Name != "" {
		m.ToolName = &msg.Name
	}

	return s.MessageRepository.Create(ctx, m)
}

func (s *store) Clear(ctx context.Context, chatID uuid.UUID) error {
	return s.MessageRepository.DeleteByChat(ctx, chatID)
}

func (s *store) Trim(ctx context.Context, chatID uuid.UUID, maxMessages int) error {
	return s.MessageRepository.DeleteByChatExceptLast(ctx, chatID, maxMessages)
}

func NewStore(chatRepository repositories.ChatRepository, messageRepository repositories.MessageRepository) Store {
	return &store{
		ChatRepository:    chatRepository,
		MessageRepository: messageRepository,
	}
}
