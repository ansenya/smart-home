package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
	RoleTool      MessageRole = "tool"
)

type MessageStatus string

const (
	StatusPending   MessageStatus = "pending"
	StatusCompleted MessageStatus = "completed"
	StatusFailed    MessageStatus = "failed"
)

type Chat struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"-"`
	Model     string         `gorm:"size:64;not null" json:"model"`
	Title     string         `gorm:"not null;default:'New Chat'" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Messages []Message `gorm:"foreignKey:ChatID" json:"-"`
}

func (Chat) TableName() string {
	return "llm_chats"
}

type Message struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ChatID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`

	Role      MessageRole `gorm:"type:varchar(20);not null" json:"role"`
	ModelName string      `gorm:"size:64;not null" json:"model_name"`

	InputTokens  int `gorm:"default:0" json:"input_tokens,omitempty"`
	OutputTokens int `gorm:"default:0" json:"output_tokens,omitempty"`

	Content string `gorm:"type:text" json:"content"`

	ToolCallID *string         `gorm:"size:255" json:"tool_call_id,omitempty"`
	ToolName   *string         `gorm:"size:255" json:"tool_name,omitempty"`
	ToolArgs   json.RawMessage `gorm:"type:jsonb" json:"tool_args,omitempty"`
	ToolResult json.RawMessage `gorm:"type:jsonb" json:"tool_result,omitempty"`

	Status MessageStatus `gorm:"type:varchar(20);not null;default:'completed'" json:"status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Message) TableName() string {
	return "llm_chat_message"
}
