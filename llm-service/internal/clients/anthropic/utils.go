package anthropic

import (
	"encoding/json"
	"llm-service/internal/clients"

	sdk "github.com/anthropics/anthropic-sdk-go"
)

// splitMessages extracts system text and converts the rest to Anthropic MessageParams.
// Anthropic does not support a "system" role in messages — it's a top-level param.
func splitMessages(msgs []clients.Message) (system []sdk.TextBlockParam, params []sdk.MessageParam) {
	for _, m := range msgs {
		switch m.Role {
		case clients.RoleSystem:
			system = append(system, sdk.TextBlockParam{Text: m.Content})

		case clients.RoleUser:
			if len(m.ToolCalls) == 0 {
				params = append(params, sdk.NewUserMessage(sdk.NewTextBlock(m.Content)))
			}

		case clients.RoleAssistant:
			blocks := []sdk.ContentBlockParamUnion{}
			if m.Content != "" {
				blocks = append(blocks, sdk.NewTextBlock(m.Content))
			}
			for _, tc := range m.ToolCalls {
				var input any
				_ = json.Unmarshal([]byte(tc.Arguments), &input)
				blocks = append(blocks, sdk.NewToolUseBlock(tc.ID, input, tc.Name))
			}
			if len(blocks) > 0 {
				params = append(params, sdk.NewAssistantMessage(blocks...))
			}

		case clients.RoleTool:
			// Tool results go as user messages with tool_result content
			params = append(params, sdk.NewUserMessage(
				sdk.NewToolResultBlock(m.ToolCallID, m.Content, false),
			))
		}
	}
	return system, params
}

func toAnthropicTools(ts []clients.Tool) []sdk.ToolUnionParam {
	out := make([]sdk.ToolUnionParam, 0, len(ts))
	for _, t := range ts {
		var props map[string]any
		var required []string

		if t.Parameters != "" {
			var schema map[string]any
			if err := json.Unmarshal([]byte(t.Parameters), &schema); err == nil {
				if p, ok := schema["properties"].(map[string]any); ok {
					props = p
				}
				if r, ok := schema["required"].([]any); ok {
					for _, v := range r {
						if s, ok := v.(string); ok {
							required = append(required, s)
						}
					}
				}
			}
		}

		tool := sdk.ToolParam{
			Name:        t.Name,
			InputSchema: sdk.ToolInputSchemaParam{Properties: props, Required: required},
		}
		if t.Description != "" {
			tool.Description = sdk.String(t.Description)
		}
		out = append(out, sdk.ToolUnionParam{OfTool: &tool})
	}
	return out
}

func parseToolCalls(content []sdk.ContentBlockUnion) []clients.ToolCall {
	var out []clients.ToolCall
	for _, block := range content {
		if block.Type == "tool_use" {
			tc := block.AsToolUse()
			out = append(out, clients.ToolCall{
				ID:        tc.ID,
				Name:      tc.Name,
				Arguments: string(tc.Input),
			})
		}
	}
	return out
}

func extractText(content []sdk.ContentBlockUnion) string {
	for _, block := range content {
		if block.Type == "text" {
			return block.AsText().Text
		}
	}
	return ""
}
