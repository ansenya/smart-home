package clients

import (
	"context"
	"llm-service/internal/clients/tools"
)

type Message struct {
	Role       string `json:"role"`
	Content    string `json:"content,omitempty"`
	ToolCallID string `json:"toolCallId,omitempty"`
	ToolCalls  any    `json:"toolCalls,omitempty"`
}

type LLMResponse struct {
	Content     string
	ToolCall    *tools.ToolCall
	ToolCallRaw any
}

type LLMClient interface {
	Generate(ctx context.Context, messages []Message, tools []tools.ToolSpec) (*LLMResponse, error)
	GenerateStream(ctx context.Context, messages []Message, tools []tools.ToolSpec, tokenChan chan<- string) (*LLMResponse, error)
}
