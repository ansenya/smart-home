package agents

import (
	"llm-service/internal/clients"
	"llm-service/internal/models"

	"github.com/google/uuid"
)

type Channel string

const (
	ChannelWeb   Channel = "web"
	ChannelVoice Channel = "voice"
)

type ChatRequest struct {
	UserID   uuid.UUID
	ChatID   uuid.UUID
	Messages []models.Message
	Model    string
	Channel  Channel
}

type ChatResponse struct {
	Output      string
	NewMessages []clients.Message
}
