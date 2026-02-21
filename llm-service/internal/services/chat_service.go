package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"llm-service/internal/clients"
	"llm-service/internal/models"
	"llm-service/internal/repositories"
	"time"

	"github.com/google/uuid"
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
	SendMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string) (*models.Message, error)
	StreamResponse(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string, tokenChan chan string) error
}
type chatService struct {
	chatRepository repositories.ChatRepository
	toolRegistry   *ToolRegistry
	llmClient      clients.LLMClient
}

func (s *chatService) CreateChat(ctx context.Context, userID uuid.UUID, model, title string) (*models.Chat, error) {
	if model == "" {
		model = "gpt-5.2"
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

	if err := s.chatRepository.CreateChat(ctx, chat); err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}
	return chat, nil
}

func (s *chatService) GetUserChats(ctx context.Context, userID uuid.UUID, limit int) ([]models.Chat, error) {
	chats, err := s.chatRepository.GetChatsByUser(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch chats by user: %w", err)
	}
	return chats, nil
}

func (s *chatService) GetChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*models.Chat, error) {
	chat, err := s.chatRepository.GetChatByID(ctx, chatID, userID)
	if err != nil {
		return nil, ErrChatNotFound
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

	if err := s.chatRepository.UpdateChat(ctx, chat); err != nil {
		return nil, fmt.Errorf("failed to update chat: %w", err)
	}

	return chat, nil
}

func (s *chatService) DeleteChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) error {
	_, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return err
	}

	if err := s.chatRepository.DeleteChat(ctx, chatID, userID); err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *chatService) GetHistory(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit int) ([]models.Message, error) {
	_, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	msgs, err := s.chatRepository.GetMessagesByChatID(ctx, chatID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return msgs, nil
}

func (s *chatService) SendMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string) (*models.Message, error) {
	if content == "" {
		return nil, ErrInvalidContent
	}

	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return nil, err
	}

	// 1. Сохраняем сообщение пользователя
	userMsg := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    userID,
		Role:      models.RoleUser,
		ModelName: "user",
		Content:   content,
		Status:    models.StatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chatRepository.CreateMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("failed to create user message: %w", err)
	}

	// 2. Запускаем цикл генерации (с поддержкой tools)
	assistantMsg, err := s.runLLMWithTools(ctx, chat, userID, false, nil)
	if err != nil {
		// Помечаем сообщение пользователя как failed
		userMsg.Status = models.StatusFailed
		s.chatRepository.UpdateMessage(ctx, userMsg)
		return nil, fmt.Errorf("llm execution failed: %w", err)
	}

	return assistantMsg, nil
}

func (s *chatService) StreamResponse(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, content string, tokenChan chan string) error {
	defer close(tokenChan)

	if content == "" {
		return ErrInvalidContent
	}

	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return err
	}

	// 1. Сохраняем сообщение пользователя
	userMsg := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		UserID:    userID,
		Role:      models.RoleUser,
		ModelName: "user",
		Content:   content,
		Status:    models.StatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chatRepository.CreateMessage(ctx, userMsg); err != nil {
		return fmt.Errorf("failed to create user message: %w", err)
	}

	// 2. Запускаем цикл генерации со стримингом
	_, err = s.runLLMWithTools(ctx, chat, userID, true, tokenChan)
	if err != nil {
		return fmt.Errorf("llm streaming failed: %w", err)
	}

	return nil
}

// runLLMWithTools - основной цикл с поддержкой инструментов
func (s *chatService) runLLMWithTools(ctx context.Context, chat *models.Chat, userID uuid.UUID, isStream bool, tokenChan chan string) (*models.Message, error) {
	// Получаем историю для контекста
	history, err := s.chatRepository.GetMessagesByChatID(ctx, chat.ID, 50)
	if err != nil {
		history = []models.Message{}
	}

	// Конвертируем в формат клиента
	messages := s.buildMessages(history)

	// Цикл генерации (максимум 5 итераций для tools)
	var lastResponse *clients.LLMResponse
	for iteration := 0; iteration < 5; iteration++ {

		// Вызов LLM
		if isStream {
			lastResponse, err = s.llmClient.GenerateStream(ctx, messages, s.toolRegistry.GetAllSpecs(), tokenChan)
		} else {
			lastResponse, err = s.llmClient.Generate(ctx, messages, s.toolRegistry.GetAllSpecs())
		}

		if err != nil {
			return nil, fmt.Errorf("llm generate failed: %w", err)
		}

		// Если есть вызов инструмента - выполняем его
		if lastResponse.ToolCall != nil {
			// Сохраняем вызов инструмента в БД
			toolMsg := &models.Message{
				ID:         uuid.New(),
				ChatID:     chat.ID,
				UserID:     userID,
				Role:       models.RoleAssistant,
				ModelName:  chat.Model,
				Content:    "",
				ToolCallID: &lastResponse.ToolCall.ID,
				ToolName:   &lastResponse.ToolCall.Name,
				ToolArgs:   lastResponse.ToolCall.Args,
				Status:     models.StatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			s.chatRepository.CreateMessage(ctx, toolMsg)

			// Находим и выполняем инструмент
			tool := s.toolRegistry.Get(lastResponse.ToolCall.Name)
			if tool == nil {
				return nil, fmt.Errorf("tool not found: %s", lastResponse.ToolCall.Name)
			}

			result, err := tool.Call(ctx, userID, lastResponse.ToolCall.Args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			// Сохраняем результат инструмента
			resultMsg := &models.Message{
				ID:         uuid.New(),
				ChatID:     chat.ID,
				UserID:     userID,
				Role:       models.RoleTool,
				ModelName:  chat.Model,
				Content:    result,
				ToolCallID: &lastResponse.ToolCall.ID,
				ToolName:   &lastResponse.ToolCall.Name,
				ToolResult: json.RawMessage(result),
				Status:     models.StatusCompleted,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			s.chatRepository.CreateMessage(ctx, resultMsg)

			// Добавляем в контекст для следующего цикла
			messages = append(messages, clients.Message{
				Role:      "assistant",
				ToolCalls: lastResponse.ToolCallRaw,
			})
			messages = append(messages, clients.Message{
				Role:       "tool",
				Content:    result,
				ToolCallID: lastResponse.ToolCall.ID,
			})

			// Продолжаем цикл для получения финального ответа
			continue
		}

		// Финальный ответ - сохраняем и возвращаем
		assistantMsg := &models.Message{
			ID:           uuid.New(),
			ChatID:       chat.ID,
			UserID:       userID,
			Role:         models.RoleAssistant,
			ModelName:    chat.Model,
			Content:      lastResponse.Content,
			InputTokens:  0, // Можно получить из ответа клиента
			OutputTokens: 0,
			Status:       models.StatusCompleted,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.chatRepository.CreateMessage(ctx, assistantMsg); err != nil {
			return nil, err
		}

		return assistantMsg, nil
	}

	return nil, fmt.Errorf("max tool iterations exceeded")
}

func (s *chatService) buildMessages(history []models.Message) []clients.Message {
	messages := make([]clients.Message, 0, len(history))
	for _, msg := range history {
		m := clients.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		}

		if msg.Role == models.RoleAssistant && msg.ToolCallID != nil {
			// TODO: конвертировать ToolArgs в ToolCalls формат OpenAI
		}

		if msg.Role == models.RoleTool && msg.ToolCallID != nil {
			m.ToolCallID = *msg.ToolCallID
		}

		messages = append(messages, m)
	}
	return messages
}

func (s *chatService) callLLM(ctx context.Context, chatID, userID uuid.UUID, messages []clients.Message) (*clients.LLMResponse, error) {
	specs := s.toolRegistry.GetAllSpecs()

	response, err := s.llmClient.Generate(ctx, messages, specs)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	if response.ToolCall != nil {
		tool := s.toolRegistry.Get(response.ToolCall.Name)
		if tool != nil {
			result, err := tool.Call(ctx, userID, response.ToolCall.Args)
			if err != nil {
				return nil, fmt.Errorf("tool call failed: %w", err)
			}
			// Здесь можно обработать результат
			_ = result
		}
	}

	return response, nil
}

func NewChatService(repos *repositories.Container, registry *ToolRegistry, openaiClient clients.LLMClient) ChatService {
	return &chatService{
		chatRepository: repos.ChatRepository,
		toolRegistry:   registry,
		llmClient:      openaiClient,
	}
}
