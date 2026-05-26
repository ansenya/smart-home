package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"llm-service/internal/agents"
	"llm-service/internal/clients"
	"llm-service/internal/config"
	"llm-service/internal/dto"
	"llm-service/internal/models"
	"llm-service/internal/repositories"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrChatNotFound     = errors.New("chat not found")
	ErrChatAccessDenied = errors.New("chat access denied")
	ErrMessageNotFound  = errors.New("message not found")
	ErrInvalidContent   = errors.New("invalid message content")
)

type ChatService interface {
	CreateChat(ctx context.Context, userID uuid.UUID, model, title string) (*models.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error)
	GetChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*models.Chat, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, title string, model string) (*models.Chat, error)
	DeleteChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error
	GetHistory(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit int) ([]models.Message, error)

	SendMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string) (*dto.MessageResponse, error)
	StreamResponse(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string, tokenChan chan string) error
	GenerateTitle(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (string, error)
}

type chatService struct {
	chatRepository    repositories.ChatRepository
	messageRepository repositories.MessageRepository
	orchestrator      agents.Orchestrator

	log *slog.Logger
}

func NewChatService(cfg *config.Container, repos *repositories.Container, orchestrator agents.Orchestrator) ChatService {
	return &chatService{
		chatRepository:    repos.ChatRepository,
		messageRepository: repos.MessageRepository,
		orchestrator:      orchestrator,

		log: cfg.Log,
	}
}

func (s *chatService) CreateChat(ctx context.Context, userID uuid.UUID, model, title string) (*models.Chat, error) {
	if model == "" {
		model = "gpt-4o"
	}
	if title == "" {
		title = "New Chat"
	}

	chat := &models.Chat{
		ID:     uuid.New(),
		UserID: userID,
		Model:  model,
		Title:  title,
	}

	if err := s.chatRepository.Create(ctx, chat); err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}
	return chat, nil
}

func (s *chatService) GetUserChats(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error) {
	chats, err := s.chatRepository.ListByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chats by user: %w", err)
	}
	return chats, nil
}

func (s *chatService) GetChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepository.GetByID(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, fmt.Errorf("get chat failed: %w", err)
	}

	if chat.UserID != userID {
		return nil, ErrChatAccessDenied
	}

	return chat, nil
}

func (s *chatService) UpdateChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, title string, model string) (*models.Chat, error) {
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		chat.Title = title
	}
	if model != "" {
		chat.Model = model
	}
	chat.UpdatedAt = time.Now()

	if err := s.chatRepository.Update(ctx, chat); err != nil {
		return nil, fmt.Errorf("failed to update chat: %w", err)
	}

	return chat, nil
}

func (s *chatService) DeleteChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	_, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return err
	}

	if err := s.chatRepository.Delete(ctx, chatID, userID); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *chatService) GetHistory(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit int) ([]models.Message, error) {
	_, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	msgs, err := s.messageRepository.ListByChat(ctx, chatID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return msgs, nil
}

func (s *chatService) SendMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string) (*dto.MessageResponse, error) {
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	history, err := s.messageRepository.ListByChat(ctx, chatID, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Persist user message
	userMsg := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		Role:      models.RoleUser,
		Content:   content,
		ModelName: chat.Model,
		Status:    models.StatusCompleted,
	}
	if err := s.messageRepository.Create(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	allMessages := append(history, *userMsg)

	resp, err := s.orchestrator.Handle(ctx, agents.ChatRequest{
		UserID:   userID,
		ChatID:   chatID,
		Messages: allMessages,
		Model:    chat.Model,
		Channel:  agents.ChannelWeb,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Persist new messages from orchestrator (tool calls + final assistant)
	if err := s.persistNewMessages(ctx, chatID, chat.Model, resp.NewMessages); err != nil {
		s.log.Error("failed to persist orchestrator messages", slog.Any("err", err))
	}

	return &dto.MessageResponse{Content: resp.Output}, nil
}

func (s *chatService) StreamResponse(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string, tokenChan chan string) error {
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return err
	}

	history, err := s.messageRepository.ListByChat(ctx, chatID, 50)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	// Persist user message
	userMsg := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		Role:      models.RoleUser,
		Content:   content,
		ModelName: chat.Model,
		Status:    models.StatusCompleted,
	}
	if err := s.messageRepository.Create(ctx, userMsg); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	allMessages := append(history, *userMsg)

	eventChan, err := s.orchestrator.Stream(ctx, agents.ChatRequest{
		UserID:   userID,
		ChatID:   chatID,
		Messages: allMessages,
		Model:    chat.Model,
		Channel:  agents.ChannelWeb,
	})
	if err != nil {
		s.log.Error("orchestrator stream failed", slog.String("model", chat.Model), slog.Any("err", err))
		return fmt.Errorf("failed to start stream: %w", err)
	}

	var fullContent strings.Builder
	var streamErr error
	for ev := range eventChan {
		if ev.Err != nil {
			streamErr = ev.Err
		}
		if ev.Done {
			break
		}
		if ev.ContentDelta != "" {
			tokenChan <- ev.ContentDelta
			fullContent.WriteString(ev.ContentDelta)
		}
	}

	if streamErr != nil {
		s.log.Error("stream recv failed", slog.String("model", chat.Model), slog.Any("err", streamErr))
		return fmt.Errorf("stream failed: %w", streamErr)
	}
	close(tokenChan)

	// Persist assistant response
	if fullContent.Len() > 0 {
		assistantMsg := &models.Message{
			ID:        uuid.New(),
			ChatID:    chatID,
			Role:      models.RoleAssistant,
			Content:   fullContent.String(),
			ModelName: chat.Model,
			Status:    models.StatusCompleted,
		}
		if err := s.messageRepository.Create(ctx, assistantMsg); err != nil {
			s.log.Error("failed to persist assistant message", slog.Any("err", err))
		}
	}

	return nil
}

func (s *chatService) persistNewMessages(ctx context.Context, chatID uuid.UUID, modelName string, msgs []clients.Message) error {
	for _, m := range msgs {
		msg := &models.Message{
			ID:        uuid.New(),
			ChatID:    chatID,
			Role:      models.MessageRole(m.Role),
			Content:   m.Content,
			ModelName: modelName,
			Status:    models.StatusCompleted,
		}

		if m.ToolCallID != "" {
			msg.ToolCallID = &m.ToolCallID
		}
		if m.Name != "" {
			msg.ToolName = &m.Name
		}
		if len(m.ToolCalls) > 0 {
			raw, _ := json.Marshal(m.ToolCalls)
			rawMsg := json.RawMessage(raw)
			msg.ToolArgs = &rawMsg
		}
		if m.Role == clients.RoleTool && m.Content != "" {
			raw := json.RawMessage(`{"result":` + jsonQuote(m.Content) + `}`)
			msg.ToolResult = &raw
		}

		if err := s.messageRepository.Create(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func jsonQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *chatService) GenerateTitle(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (string, error) {
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return "", err
	}

	// Find first user message
	msgs, err := s.messageRepository.ListByChat(ctx, chatID, 5)
	if err != nil {
		return "", err
	}
	var firstUserContent string
	for _, m := range msgs {
		if m.Role == models.RoleUser {
			firstUserContent = m.Content
			break
		}
	}
	if firstUserContent == "" {
		return "", fmt.Errorf("no user messages found")
	}

	// Trim long messages
	if len(firstUserContent) > 400 {
		firstUserContent = firstUserContent[:400]
	}

	resp, err := s.orchestrator.Handle(ctx, agents.ChatRequest{
		UserID: userID,
		ChatID: chatID,
		Model:  chat.Model,
		Messages: []models.Message{
			{
				Role:      models.RoleUser,
				Content:   "Generate a short title (3-6 words) for a chat that starts with this message. Reply with ONLY the title, no quotes, no punctuation at the end:\n\n" + firstUserContent,
				ModelName: chat.Model,
			},
		},
	})
	if err != nil {
		return "", err
	}

	title := strings.TrimSpace(resp.Output)
	// Strip surrounding quotes if model adds them
	title = strings.Trim(title, `"'`)
	if len(title) > 100 {
		title = title[:100]
	}

	chat.Title = title
	chat.UpdatedAt = time.Now()
	if err := s.chatRepository.Update(ctx, chat); err != nil {
		return "", err
	}

	return title, nil
}
