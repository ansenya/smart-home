package openai

import (
	"context"
	"io"
	"llm-service/internal/clients"

	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/responses"
)

type providerStream struct {
	sdk      *ssestream.Stream[responses.ResponseStreamEventUnion]
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

		ev := s.sdk.Current()

		switch v := ev.AsAny().(type) {
		case responses.ResponseTextDeltaEvent:
			if v.Delta != "" {
				return &clients.StreamEvent{ContentDelta: v.Delta}, nil
			}

		case responses.ResponseFunctionCallArgumentsDeltaEvent:
			return &clients.StreamEvent{
				ToolCallDelta: &clients.ToolCall{
					ID:        v.ItemID,
					Arguments: v.Delta,
				},
			}, nil

		case responses.ResponseFunctionCallArgumentsDoneEvent:
			return &clients.StreamEvent{
				ToolCallDelta: &clients.ToolCall{
					ID:        v.ItemID,
					Name:      v.Name,
					Arguments: v.Arguments,
				},
			}, nil

		case responses.ResponseCompletedEvent:
			s.finished = true
			return &clients.StreamEvent{Done: true}, nil

		default:
			continue
		}
	}
}

func (s *providerStream) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return nil
}
