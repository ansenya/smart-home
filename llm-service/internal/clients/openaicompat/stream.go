package openaicompat

import (
	"context"
	"io"
	"llm-service/internal/clients"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

type providerStream struct {
	sdk      *ssestream.Stream[openai.ChatCompletionChunk]
	cancel   context.CancelFunc
	finished bool
}

func (s *providerStream) Recv(ctx context.Context) (*clients.StreamEvent, error) {
	if s.finished {
		return nil, io.EOF
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if !s.sdk.Next() {
			if err := s.sdk.Err(); err != nil {
				return nil, err
			}
			s.finished = true
			return &clients.StreamEvent{Done: true}, nil
		}

		chunk := s.sdk.Current()
		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]

		if choice.Delta.Content != "" {
			return &clients.StreamEvent{ContentDelta: choice.Delta.Content}, nil
		}

		if choice.FinishReason != "" {
			s.finished = true
			return &clients.StreamEvent{Done: true}, nil
		}
	}
}

func (s *providerStream) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}
