package runtime

import (
	"context"
	"fmt"

	"llm-service/internal/clients"
)

type ToolLoop interface {
	Run(ctx context.Context, provider clients.Provider, req *clients.Request) (*clients.Response, error)
	RunStream(ctx context.Context, provider clients.Provider, req *clients.Request) (clients.Stream, error)
}

type ToolExecutor interface {
	Call(ctx context.Context, name string, args string) (string, error)
}

type ToolFunc func(ctx context.Context, args string) (string, error)

type ToolMap map[string]ToolFunc

func (m ToolMap) Call(ctx context.Context, name string, args string) (string, error) {
	fn, ok := m[name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return fn(ctx, args)
}

type loop struct {
	exec     ToolExecutor
	maxIters int
}

func NewToolLoop(exec ToolExecutor, maxIters int) ToolLoop {
	if maxIters <= 0 {
		maxIters = 8
	}
	return &loop{exec: exec, maxIters: maxIters}
}

func (l *loop) Run(ctx context.Context, provider clients.Provider, req *clients.Request) (*clients.Response, error) {
	if l.exec == nil {
		return nil, fmt.Errorf("tool executor is nil")
	}

	for i := 0; i < l.maxIters; i++ {
		resp, err := provider.Generate(ctx, req)
		if err != nil {
			return nil, err
		}

		if len(resp.ToolCalls) == 0 {
			return resp, nil
		}

		req.Messages = append(req.Messages, clients.Message{
			Role:      clients.RoleAssistant,
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		for _, tc := range resp.ToolCalls {
			out, err := l.exec.Call(ctx, tc.Name, tc.Arguments)
			if err != nil {
				return nil, fmt.Errorf("tool %s failed: %w", tc.Name, err)
			}

			req.Messages = append(req.Messages, clients.Message{
				Role:       clients.RoleTool,
				Name:       tc.Name,
				ToolCallID: tc.ID,
				Content:    out,
			})
		}
	}

	return nil, fmt.Errorf("tool loop: max iterations reached")
}

// RunStream runs the tool loop using Generate for intermediate steps and
// returns a Stream for the final (non-tool-call) response.
func (l *loop) RunStream(ctx context.Context, provider clients.Provider, req *clients.Request) (clients.Stream, error) {
	if l.exec == nil {
		return nil, fmt.Errorf("tool executor is nil")
	}

	for i := 0; i < l.maxIters; i++ {
		resp, err := provider.Generate(ctx, req)
		if err != nil {
			return nil, err
		}

		if len(resp.ToolCalls) == 0 {
			// Final step — stream it
			streamReq := *req
			streamReq.Stream = true
			return provider.Stream(ctx, &streamReq)
		}

		req.Messages = append(req.Messages, clients.Message{
			Role:      clients.RoleAssistant,
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		for _, tc := range resp.ToolCalls {
			out, err := l.exec.Call(ctx, tc.Name, tc.Arguments)
			if err != nil {
				return nil, fmt.Errorf("tool %s failed: %w", tc.Name, err)
			}

			req.Messages = append(req.Messages, clients.Message{
				Role:       clients.RoleTool,
				Name:       tc.Name,
				ToolCallID: tc.ID,
				Content:    out,
			})
		}
	}

	return nil, fmt.Errorf("tool loop: max iterations reached")
}
