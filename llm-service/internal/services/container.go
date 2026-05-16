package services

import (
	"llm-service/internal/agents"
	"llm-service/internal/config"
	"llm-service/internal/repositories"
)

type Container struct {
	ChatService ChatService
}

func NewContainer(cfg *config.Container, repos *repositories.Container, orchestrator agents.Orchestrator) (*Container, error) {
	return &Container{
		ChatService: NewChatService(cfg, repos, orchestrator),
	}, nil
}
