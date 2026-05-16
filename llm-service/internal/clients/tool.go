package clients

type Tool struct {
	Name        string
	Description string
	Parameters  string
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}
