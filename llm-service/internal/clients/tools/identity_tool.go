package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"llm-service/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdentityTool interface {
	Tool
}

type identityTool struct {
	accountRepository repositories.UsersRepository
}

func (t *identityTool) Name() string {
	return "get_users_identity"
}

func (t *identityTool) Description() string {
	return "Get user's identity information: name"
}

func (t *identityTool) JSONSchema() any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{""},
	}
}

func (t *identityTool) Call(ctx context.Context, userID uuid.UUID, _ json.RawMessage) (string, error) {
	user, err := t.accountRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", fmt.Errorf("failed to get user by id: %v", err)
	}

	resp := map[string]any{
		"name": user.Name,
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func NewIdentityTool(accountRepository repositories.UsersRepository) IdentityTool {
	return &identityTool{
		accountRepository: accountRepository,
	}
}
