package services

import (
	"fmt"
	"llm-service/internal/clients"
	"llm-service/internal/config"
	"llm-service/internal/repositories"
)

type Container struct {
	ChatService ChatService
}

func NewContainer(cfg *config.Container, repos *repositories.Container) (*Container, error) {
	toolRegistry := NewToolRegistry()
	openaiClient, err := clients.NewOpenAIClient(cfg.OpenaiConfig.ApiKey, cfg.OpenaiConfig.ProxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create open ai client: %w", err)
	}
	return &Container{
		ChatService: NewChatService(repos, toolRegistry),
	}, nil
}
