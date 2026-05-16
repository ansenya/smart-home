package clients

type Response struct {
	Content       string
	ToolCalls     []ToolCall
	Usage         *Usage
	NewResponseID string
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}
