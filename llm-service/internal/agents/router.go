package agents

import (
	"context"
	"fmt"
	"llm-service/internal/clients"
	"strings"
)

type ModelName = string

type ProviderRouter interface {
	Register(p clients.Provider, models ...string)
	RegisterPrefix(p clients.Provider, prefix string)
	Resolve(ctx context.Context, model string) (clients.Provider, error)
}

type prefixEntry struct {
	prefix   string
	provider clients.Provider
}

type providerRegistry struct {
	exact    map[string]clients.Provider
	prefixes []prefixEntry
}

func NewProviderRegistry() ProviderRouter {
	return &providerRegistry{
		exact: make(map[string]clients.Provider),
	}
}

func (r *providerRegistry) Register(p clients.Provider, models ...string) {
	for _, m := range models {
		r.exact[m] = p
	}
}

func (r *providerRegistry) RegisterPrefix(p clients.Provider, prefix string) {
	r.prefixes = append(r.prefixes, prefixEntry{prefix: prefix, provider: p})
}

func (r *providerRegistry) Resolve(_ context.Context, model string) (clients.Provider, error) {
	if p, ok := r.exact[model]; ok {
		return p, nil
	}
	for _, pe := range r.prefixes {
		if strings.HasPrefix(model, pe.prefix) {
			return pe.provider, nil
		}
	}
	return nil, fmt.Errorf("no provider found for model %q", model)
}
