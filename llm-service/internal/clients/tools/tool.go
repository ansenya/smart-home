package tools

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Tool interface {
	Name() string
	Description() string
	JSONSchema() any
	Call(ctx context.Context, userID uuid.UUID, args json.RawMessage) (string, error)
}

type ToolSpec struct {
	Name        string
	Description string
	Schema      any
}

type ToolCall struct {
	ID   string
	Name string
	Args json.RawMessage
}
