package anthropic

import (
	"context"
	"io"
	"llm-service/internal/clients"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

type providerStream struct {
	sdk      *ssestream.Stream[sdk.MessageStreamEventUnion]
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
		switch ev.Type {
		case "content_block_delta":
			delta := ev.AsContentBlockDelta()
			raw := delta.Delta
			if raw.Type == "text_delta" {
				text := raw.AsTextDelta()
				if text.Text != "" {
					return &clients.StreamEvent{ContentDelta: text.Text}, nil
				}
			}
		case "message_stop":
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
