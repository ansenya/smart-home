package clients

type Request struct {
	Model    string
	Messages []Message
	Tools    []Tool

	Temperature float32
	MaxTokens   int

	Stream bool

	PreviousResponseID string
}
