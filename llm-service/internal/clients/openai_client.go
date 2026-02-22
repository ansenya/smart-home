package clients

import (
	"context"
	"errors"
	"fmt"
	"io"
	"llm-service/internal/clients/tools"
	"net/http"
	"net/url"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type OpenaiClient struct {
	api *openai.Client
}

func (c *OpenaiClient) GenerateStream(ctx context.Context, messages []Message, specs []tools.ToolSpec, tokenChan chan<- string) (*LLMResponse, error) {
	req := openai.ChatCompletionRequest{
		Model:               "gpt-5.2",
		Temperature:         0.3,
		MaxCompletionTokens: 1500,
		Stream:              true, // Включаем стриминг
		Messages: func() []openai.ChatCompletionMessage {
			out := make([]openai.ChatCompletionMessage, len(messages))
			for i, m := range messages {
				msg := openai.ChatCompletionMessage{
					Role:    m.Role,
					Content: m.Content,
				}

				if msg.Content == "" {
					msg.Content = ""
				}

				if m.ToolCalls != nil {
					msg.ToolCalls = m.ToolCalls.([]openai.ToolCall)
				}

				if m.ToolCallID != "" {
					msg.ToolCallID = m.ToolCallID
				}

				out[i] = msg
			}
			return out
		}(),
		Tools: func() []openai.Tool {
			out := make([]openai.Tool, len(specs))
			for i, t := range specs {
				out[i] = openai.Tool{
					Type: openai.ToolTypeFunction,
					Function: &openai.FunctionDefinition{
						Name:        t.Name,
						Description: t.Description,
						Parameters:  t.Schema,
					},
				}
			}
			return out
		}(),
	}

	// Создаем стрим
	stream, err := c.api.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	var fullContent strings.Builder
	var toolCalls []openai.ToolCall
	var hasToolCall bool

	// Читаем поток токенов
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("stream error: %w", err)
		}

		if len(response.Choices) == 0 {
			continue
		}

		delta := response.Choices[0].Delta

		// Отправляем токен в канал
		if delta.Content != "" {
			select {
			case tokenChan <- delta.Content:
				fullContent.WriteString(delta.Content)
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		// Собираем tool_calls (они могут приходить частями)
		if len(delta.ToolCalls) > 0 {
			hasToolCall = true
			for i, tc := range delta.ToolCalls {
				if i >= len(toolCalls) {
					toolCalls = append(toolCalls, tc)
				} else {
					// Объединяем аргументы если они приходят частями
					toolCalls[i].Function.Arguments += tc.Function.Arguments
					if tc.Function.Name != "" {
						toolCalls[i].Function.Name = tc.Function.Name
					}
					if tc.ID != "" {
						toolCalls[i].ID = tc.ID
					}
				}
			}
		}
	}

	// Если был вызов инструмента - возвращаем его
	if hasToolCall && len(toolCalls) > 0 {
		call := toolCalls[0]
		return &LLMResponse{
			ToolCall: &tools.ToolCall{
				ID:   call.ID,
				Name: call.Function.Name,
				Args: []byte(call.Function.Arguments),
			},
			ToolCallRaw: toolCalls,
		}, nil
	}

	// Иначе возвращаем текстовый ответ
	return &LLMResponse{
		Content: fullContent.String(),
	}, nil
}

func NewOpenAIClient(apiKey string, proxyURL string) (LLMClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be blank")
	}

	config := openai.DefaultConfig(apiKey)

	if proxyURL != "" {
		proxyUrl, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy URL: %w", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	return &OpenaiClient{api: openai.NewClientWithConfig(config)}, nil
}

func (c *OpenaiClient) Generate(ctx context.Context, messages []Message, specs []tools.ToolSpec) (*LLMResponse, error) {
	req := openai.ChatCompletionRequest{
		Model:               "gpt-5.2",
		MaxCompletionTokens: 1500,
		Messages: func() []openai.ChatCompletionMessage {
			out := make([]openai.ChatCompletionMessage, len(messages))
			for i, m := range messages {
				msg := openai.ChatCompletionMessage{
					Role:    m.Role,
					Content: m.Content,
				}

				// OpenAI требует string, не nil
				if msg.Content == "" {
					msg.Content = ""
				}

				// assistant → tool_calls
				if m.ToolCalls != nil {
					msg.ToolCalls = m.ToolCalls.([]openai.ToolCall)
				}

				// tool → tool_call_id
				if m.ToolCallID != "" {
					msg.ToolCallID = m.ToolCallID
				}

				out[i] = msg
			}
			return out
		}(),
		Tools: func() []openai.Tool {
			out := make([]openai.Tool, len(specs))
			for i, t := range specs {
				out[i] = openai.Tool{
					Type: openai.ToolTypeFunction,
					Function: &openai.FunctionDefinition{
						Name:        t.Name,
						Description: t.Description,
						Parameters:  t.Schema,
					},
				}
			}
			return out
		}(),
	}

	resp, err := c.api.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) > 0 {
		call := msg.ToolCalls[0]

		return &LLMResponse{
			ToolCall: &tools.ToolCall{
				ID:   call.ID,
				Name: call.Function.Name,
				Args: []byte(call.Function.Arguments),
			},
			ToolCallRaw: msg.ToolCalls,
		}, nil
	}

	return &LLMResponse{
		Content: msg.Content,
	}, nil
}
