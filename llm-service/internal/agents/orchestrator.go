package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"llm-service/internal/agents/runtime"
	"llm-service/internal/clients"
	"llm-service/internal/ctxkeys"
	"llm-service/internal/models"
)

type Orchestrator interface {
	Handle(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	Stream(ctx context.Context, req ChatRequest) (<-chan clients.StreamEvent, error)
}

type orchestrator struct {
	router   ProviderRouter
	toolLoop runtime.ToolLoop
	tools    []clients.Tool
}

func NewOrchestrator(router ProviderRouter, toolLoop runtime.ToolLoop, tools []clients.Tool) Orchestrator {
	return &orchestrator{router: router, toolLoop: toolLoop, tools: tools}
}

func (o *orchestrator) Handle(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	ctx = context.WithValue(ctx, ctxkeys.UserID, req.UserID)

	provider, err := o.router.Resolve(ctx, req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model %s: %w", req.Model, err)
	}

	clientReq := &clients.Request{
		Model:    req.Model,
		Messages: convertMessages(req.Messages),
		Tools:    o.tools,
	}

	startLen := len(clientReq.Messages)

	var resp *clients.Response
	if o.toolLoop != nil && len(o.tools) > 0 {
		resp, err = o.toolLoop.Run(ctx, provider, clientReq)
	} else {
		resp, err = provider.Generate(ctx, clientReq)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	// Collect messages appended during tool loop (assistant + tool result turns)
	newMessages := clientReq.Messages[startLen:]
	// Append final assistant response
	if resp.Content != "" {
		newMessages = append(newMessages, clients.Message{
			Role:    clients.RoleAssistant,
			Content: resp.Content,
		})
	}

	return &ChatResponse{Output: resp.Content, NewMessages: newMessages}, nil
}

func (o *orchestrator) Stream(ctx context.Context, req ChatRequest) (<-chan clients.StreamEvent, error) {
	ctx = context.WithValue(ctx, ctxkeys.UserID, req.UserID)

	provider, err := o.router.Resolve(ctx, req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model %s: %w", req.Model, err)
	}

	clientReq := &clients.Request{
		Model:    req.Model,
		Messages: convertMessages(req.Messages),
		Tools:    o.tools,
		Stream:   true,
	}

	var stream clients.Stream
	if o.toolLoop != nil && len(o.tools) > 0 {
		stream, err = o.toolLoop.RunStream(ctx, provider, clientReq)
	} else {
		stream, err = provider.Stream(ctx, clientReq)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to start stream: %w", err)
	}

	ch := make(chan clients.StreamEvent, 32)
	go func() {
		defer close(ch)
		defer stream.Close()
		for {
			ev, err := stream.Recv(ctx)
			if err != nil {
				if err != io.EOF {
					ch <- clients.StreamEvent{Done: true, Err: err}
				} else {
					ch <- clients.StreamEvent{Done: true}
				}
				return
			}
			ch <- *ev
			if ev.Done {
				return
			}
		}
	}()

	return ch, nil
}

func convertMessages(msgs []models.Message) []clients.Message {
	out := make([]clients.Message, 0, len(msgs))
	for _, m := range msgs {
		cm := clients.Message{
			Role:    clients.Role(m.Role),
			Content: m.Content,
		}
		if m.ToolCallID != nil {
			cm.ToolCallID = *m.ToolCallID
		}
		if m.ToolName != nil {
			cm.Name = *m.ToolName
		}
		// Reconstruct tool_calls for assistant messages from stored ToolArgs.
		// Without this, history replayed to the LLM would contain tool-result messages
		// with no preceding assistant tool_call, causing a 400 from OpenAI/Anthropic.
		if m.Role == models.RoleAssistant && m.ToolArgs != nil {
			var calls []clients.ToolCall
			if err := json.Unmarshal(*m.ToolArgs, &calls); err == nil && len(calls) > 0 {
				cm.ToolCalls = calls
			}
		}
		out = append(out, cm)
	}
	return out
}
