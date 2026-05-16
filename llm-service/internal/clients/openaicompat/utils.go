package openaicompat

import (
	"encoding/json"
	"llm-service/internal/clients"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

func toMessages(msgs []clients.Message) []openai.ChatCompletionMessageParamUnion {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(msgs))
	for _, m := range msgs {
		switch m.Role {
		case clients.RoleSystem:
			out = append(out, openai.SystemMessage(m.Content))
		case clients.RoleUser:
			out = append(out, openai.UserMessage(m.Content))
		case clients.RoleAssistant:
			asst := openai.ChatCompletionAssistantMessageParam{}
			asst.Content.OfString = openai.String(m.Content)
			if len(m.ToolCalls) > 0 {
				asst.ToolCalls = toToolCallParams(m.ToolCalls)
			}
			out = append(out, openai.ChatCompletionMessageParamUnion{OfAssistant: &asst})
		case clients.RoleTool:
			out = append(out, openai.ToolMessage(m.Content, m.ToolCallID))
		}
	}
	return out
}

func toToolCallParams(tcs []clients.ToolCall) []openai.ChatCompletionMessageToolCallUnionParam {
	out := make([]openai.ChatCompletionMessageToolCallUnionParam, 0, len(tcs))
	for _, tc := range tcs {
		ft := openai.ChatCompletionMessageFunctionToolCallParam{
			ID: tc.ID,
			Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
				Name:      tc.Name,
				Arguments: tc.Arguments,
			},
		}
		out = append(out, openai.ChatCompletionMessageToolCallUnionParam{
			OfFunction: &ft,
		})
	}
	return out
}

func toTools(ts []clients.Tool) []openai.ChatCompletionToolUnionParam {
	out := make([]openai.ChatCompletionToolUnionParam, 0, len(ts))
	for _, t := range ts {
		var params map[string]any
		if t.Parameters != "" {
			_ = json.Unmarshal([]byte(t.Parameters), &params)
		}
		fn := shared.FunctionDefinitionParam{
			Name:       t.Name,
			Parameters: shared.FunctionParameters(params),
		}
		if t.Description != "" {
			fn.Description = openai.String(t.Description)
		}
		out = append(out, openai.ChatCompletionToolUnionParam{
			OfFunction: &openai.ChatCompletionFunctionToolParam{
				Function: fn,
			},
		})
	}
	return out
}

func parseToolCalls(tcs []openai.ChatCompletionMessageToolCallUnion) []clients.ToolCall {
	out := make([]clients.ToolCall, 0, len(tcs))
	for _, tc := range tcs {
		fn := tc.AsFunction()
		out = append(out, clients.ToolCall{
			ID:        fn.ID,
			Name:      fn.Function.Name,
			Arguments: fn.Function.Arguments,
		})
	}
	return out
}
