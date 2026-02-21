package dto

import (
	"time"

	"github.com/google/uuid"
)

// requests

type CreateChatRequest struct {
	Model string `json:"model,omitempty"`
	Title string `json:"title,omitempty"`
}

type UpdateChatRequest struct {
	Title string `json:"title,omitempty"`
	Model string `json:"model,omitempty"`
}

type SendMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=10000"`
}

// responses

type ChatResponse struct {
	ID        uuid.UUID `json:"id"`
	Model     string    `json:"model"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatListResponse struct {
	Chats []ChatResponse `json:"chats"`
}

type MessageResponse struct {
	ID           uuid.UUID      `json:"id"`
	Role         string         `json:"role"`
	Content      string         `json:"content"`
	ModelName    string         `json:"model_name,omitempty"`
	InputTokens  int            `json:"input_tokens,omitempty"`
	OutputTokens int            `json:"output_tokens,omitempty"`
	ToolCallID   *string        `json:"tool_call_id,omitempty"`
	ToolName     *string        `json:"tool_name,omitempty"`
	ToolArgs     map[string]any `json:"tool_args,omitempty"`
	ToolResult   map[string]any `json:"tool_result,omitempty"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}

type MessageListResponse struct {
	Messages []MessageResponse `json:"messages"`
	HasMore  bool              `json:"has_more"`
}

type StreamChunk struct {
	Token string `json:"token"`
	Done  bool   `json:"done"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
