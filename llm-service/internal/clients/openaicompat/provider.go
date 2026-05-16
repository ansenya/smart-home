package openaicompat

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"llm-service/internal/clients"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
)

type provider struct {
	client *openai.Client
	name   string
	apiKey string
}

func (p *provider) Name() string {
	return p.name
}

func (p *provider) Capabilities() clients.Capability {
	return clients.Capability{
		Tools:      true,
		Streaming:  true,
		Vision:     false,
		JSONMode:   true,
		Embeddings: false,
	}
}

func (p *provider) Generate(ctx context.Context, req *clients.Request) (*clients.Response, error) {
	params := openai.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: toMessages(req.Messages),
	}
	if len(req.Tools) > 0 {
		params.Tools = toTools(req.Tools)
	}
	if req.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(int64(req.MaxTokens))
	}
	if req.Temperature > 0 {
		params.Temperature = param.NewOpt(float64(req.Temperature))
	}

	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("%s provider generate: %w", p.name, err)
	}

	if len(resp.Choices) == 0 {
		return &clients.Response{}, nil
	}

	msg := resp.Choices[0].Message
	return &clients.Response{
		Content:   msg.Content,
		ToolCalls: parseToolCalls(msg.ToolCalls),
		Usage: &clients.Usage{
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		},
	}, nil
}

func (p *provider) Stream(ctx context.Context, req *clients.Request) (clients.Stream, error) {
	params := openai.ChatCompletionNewParams{
		Model:    req.Model,
		Messages: toMessages(req.Messages),
	}
	if len(req.Tools) > 0 {
		params.Tools = toTools(req.Tools)
	}
	if req.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(int64(req.MaxTokens))
	}

	streamCtx, cancel := context.WithCancel(ctx)
	sdk := p.client.Chat.Completions.NewStreaming(streamCtx, params)
	return &providerStream{sdk: sdk, cancel: cancel}, nil
}

func (p *provider) WithAPIKey(key string) clients.Provider {
	opts := []option.RequestOption{option.WithAPIKey(key)}
	// re-use same base URL by extracting from original client via a new client
	newClient := openai.NewClient(opts...)
	return &provider{client: &newClient, name: p.name, apiKey: key}
}

func New(name, baseURL, apiKey, proxyURL string) (clients.Provider, error) {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		opts = append(opts, option.WithHTTPClient(&http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(parsed)},
		}))
	}
	c := openai.NewClient(opts...)
	return &provider{client: &c, name: name, apiKey: apiKey}, nil
}
