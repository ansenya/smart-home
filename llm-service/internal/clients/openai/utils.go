package openai

import (
	"encoding/json"
	"llm-service/internal/clients"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func toInputItems(hist []clients.Message) []responses.ResponseInputItemUnionParam {
	items := make([]responses.ResponseInputItemUnionParam, 0, len(hist))

	for _, m := range hist {
		if m.Role == clients.RoleAssistant {
			continue
		}

		items = append(items, responses.ResponseInputItemUnionParam{
			OfInputMessage: &responses.ResponseInputItemMessageParam{
				Role: string(m.Role),
				Content: []responses.ResponseInputContentUnionParam{
					{
						OfInputText: &responses.ResponseInputTextParam{
							Text: m.Content,
						},
					},
				},
			},
		})
	}

	return items
}

func toOpenAITools(ts []clients.Tool) []responses.ToolUnionParam {
	out := make([]responses.ToolUnionParam, 0, len(ts))
	for _, t := range ts {
		var params map[string]any
		if t.Parameters != "" {
			if err := json.Unmarshal([]byte(t.Parameters), &params); err != nil {
				// фоллбек на пустой объект-схему
				params = map[string]any{"type": "object", "properties": map[string]any{}}
			}
		}

		ft := responses.FunctionToolParam{
			Name: t.Name,
		}
		if t.Description != "" {
			ft.Description = openai.String(t.Description)
		}
		if params != nil {
			ft.Parameters = params
		}

		out = append(out, responses.ToolUnionParam{
			OfFunction: &ft,
		})
	}
	return out

}

func parseToolCalls(resp *responses.Response) []clients.ToolCall {
	var calls []clients.ToolCall
	for _, item := range resp.Output {
		if item.Type == "function_call" {
			fc := item.AsFunctionCall()
			calls = append(calls, clients.ToolCall{
				ID:        fc.CallID,
				Name:      fc.Name,
				Arguments: fc.Arguments,
			})
		}
	}
	return calls
}

func mapUsage(resp *responses.Response) *clients.Usage {
	return &clients.Usage{
		PromptTokens:     int(resp.Usage.InputTokens),
		CompletionTokens: int(resp.Usage.OutputTokens),
		TotalTokens:      int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
	}
}
