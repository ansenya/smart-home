package openai

import (
	"context"
	"fmt"
	"llm-service/internal/clients"
	"net/http"
	"net/url"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

type provider struct {
	client *openai.Client
}

func (p *provider) Name() string {
	return "openai"
}

func (p *provider) Capabilities() clients.Capability {
	return clients.Capability{
		Tools:      true,
		Streaming:  true,
		Vision:     true,
		JSONMode:   true,
		Embeddings: true,
	}
}

func (p *provider) Generate(ctx context.Context, req *clients.Request) (*clients.Response, error) {
	if req.Model == "" {
		req.Model = openai.ChatModelGPT5_2
	}
	if req.MaxTokens < 16 {
		req.MaxTokens = 16
	}

	var previousResponseID param.Opt[string]
	if req.PreviousResponseID != "" {
		previousResponseID = openai.String(req.PreviousResponseID)
	}
	resp, err := p.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: req.Model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toInputItems(req.Messages),
		},
		Tools:              toOpenAITools(req.Tools),
		Temperature:        openai.Float(float64(req.Temperature)),
		MaxOutputTokens:    openai.Int(int64(req.MaxTokens)),
		PreviousResponseID: previousResponseID,
	})
	if err != nil {
		return nil, fmt.Errorf("openai provider failed to generate: %w", err)
	}

	return &clients.Response{
		Content:   resp.OutputText(),
		ToolCalls: parseToolCalls(resp),
		Usage:     mapUsage(resp),
	}, nil
}

func (p *provider) Stream(ctx context.Context, req *clients.Request) (clients.Stream, error) {
	if req.Model == "" {
		req.Model = openai.ChatModelGPT5_2
	}
	if req.MaxTokens < 16 {
		req.MaxTokens = 16
	}

	streamCtx, cancel := context.WithCancel(ctx)
	sdk := p.client.Responses.NewStreaming(streamCtx, responses.ResponseNewParams{
		Model: req.Model,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: toInputItems(req.Messages),
		},
		Tools:           toOpenAITools(req.Tools),
		Temperature:     openai.Float(float64(req.Temperature)),
		MaxOutputTokens: openai.Int(int64(req.MaxTokens)),
	})

	return &providerStream{sdk: sdk, cancel: cancel}, nil
}

func New(apiKey string, proxyURL string) (clients.Provider, error) {
	options := []option.RequestOption{option.WithAPIKey(apiKey)}
	if proxyURL != "" {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parsedURL),
			},
		}
		options = append(options, option.WithHTTPClient(httpClient))
	}
	client := openai.NewClient(options...)
	return &provider{
		client: &client,
	}, nil
}
