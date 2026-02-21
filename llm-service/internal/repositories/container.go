package repositories

import (
	"llm-service/internal/infra/db"
)

type Container struct {
	UsersRepository   UsersRepository
	ChatRepository    ChatRepository
	SessionRepository SessionRepository
}

func NewContainer(adapter *db.Client) *Container {
	return &Container{
		UsersRepository:   NewUsersRepository(adapter.DB),
		ChatRepository:    NewChatRepository(adapter.DB),
		SessionRepository: NewSessionRepository(adapter.DB),
	}
}
