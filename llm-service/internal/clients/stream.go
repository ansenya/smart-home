package clients

import "context"

type Stream interface {
	Recv(ctx context.Context) (*StreamEvent, error)
	Close() error
}

type StreamEvent struct {
	ContentDelta  string
	ToolCallDelta *ToolCall
	Done          bool
	Err           error
}
