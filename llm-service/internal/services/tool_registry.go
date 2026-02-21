package services

import "llm-service/internal/clients/tools"

type ToolRegistry struct {
	tools map[string]tools.Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]tools.Tool),
	}
}

func (r *ToolRegistry) Register(tool tools.Tool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) tools.Tool {
	return r.tools[name]
}

func (r *ToolRegistry) GetAllSpecs() []tools.ToolSpec {
	specs := make([]tools.ToolSpec, 0, len(r.tools))
	for _, tool := range r.tools {
		specs = append(specs, tools.ToolSpec{
			Name:        tool.Name(),
			Description: tool.Description(),
			Schema:      tool.JSONSchema(),
		})
	}
	return specs
}
