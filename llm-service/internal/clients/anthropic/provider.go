package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"llm-service/internal/clients"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const defaultMaxTokens = 8192

type provider struct {
	client *sdk.Client
}

func (p *provider) Name() string {
	return "anthropic"
}

func (p *provider) Capabilities() clients.Capability {
	return clients.Capability{
		Tools:      true,
		Streaming:  true,
		Vision:     true,
		JSONMode:   false,
		Embeddings: false,
	}
}

func (p *provider) Generate(ctx context.Context, req *clients.Request) (*clients.Response, error) {
	system, messages := splitMessages(req.Messages)

	maxTokens := int64(req.MaxTokens)
	if maxTokens < 1 {
		maxTokens = defaultMaxTokens
	}

	params := sdk.MessageNewParams{
		Model:    sdk.Model(req.Model),
		Messages: messages,
		MaxTokens: maxTokens,
	}
	if len(system) > 0 {
		params.System = system
	}
	if len(req.Tools) > 0 {
		params.Tools = toAnthropicTools(req.Tools)
	}

	resp, err := p.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic generate: %w", err)
	}

	return &clients.Response{
		Content:   extractText(resp.Content),
		ToolCalls: parseToolCalls(resp.Content),
		Usage: &clients.Usage{
			PromptTokens:     int(resp.Usage.InputTokens),
			CompletionTokens: int(resp.Usage.OutputTokens),
			TotalTokens:      int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
		},
	}, nil
}

func (p *provider) Stream(ctx context.Context, req *clients.Request) (clients.Stream, error) {
	system, messages := splitMessages(req.Messages)

	maxTokens := int64(req.MaxTokens)
	if maxTokens < 1 {
		maxTokens = defaultMaxTokens
	}

	params := sdk.MessageNewParams{
		Model:    sdk.Model(req.Model),
		Messages: messages,
		MaxTokens: maxTokens,
	}
	if len(system) > 0 {
		params.System = system
	}
	if len(req.Tools) > 0 {
		params.Tools = toAnthropicTools(req.Tools)
	}

	streamCtx, cancel := context.WithCancel(ctx)
	stream := p.client.Messages.NewStreaming(streamCtx, params)
	return &providerStream{sdk: stream, cancel: cancel}, nil
}

func (p *provider) WithAPIKey(key string) clients.Provider {
	c := sdk.NewClient(option.WithAPIKey(key))
	return &provider{client: &c}
}

func New(apiKey, proxyURL string) (clients.Provider, error) {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		opts = append(opts, option.WithHTTPClient(&http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(parsed)},
		}))
	}
	c := sdk.NewClient(opts...)
	return &provider{client: &c}, nil
}
