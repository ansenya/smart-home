package repositories

import (
	"llm-service/internal/infra/db"
)

type Container struct {
	UsersRepository   UsersRepository
	ChatRepository    ChatRepository
	MessageRepository MessageRepository
	SessionRepository SessionRepository
}

func NewContainer(adapter *db.Client) *Container {
	return &Container{
		UsersRepository:   NewUsersRepository(adapter.DB),
		ChatRepository:    NewChatRepository(adapter.DB),
		MessageRepository: NewMessageRepository(adapter.DB),
		SessionRepository: NewSessionRepository(adapter.DB),
	}
}
