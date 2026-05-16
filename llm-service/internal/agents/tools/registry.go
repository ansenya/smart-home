package tools

import "context"

type Handler func(ctx context.Context, args []byte) (string, error)

type Registry interface {
	Register(name string, handler Handler)
	Get(name string) (Handler, bool)
	List() []string
}

type registry struct {
	m map[string]Handler
}

func NewRegistry() Registry {
	return &registry{m: make(map[string]Handler)}
}

func (r *registry) Register(name string, handler Handler) {
	if r.m == nil {
		r.m = make(map[string]Handler)
	}
	r.m[name] = handler
}

func (r *registry) Get(name string) (Handler, bool) {
	h, ok := r.m[name]
	return h, ok
}

func (r *registry) List() []string {
	out := make([]string, 0, len(r.m))
	for k := range r.m {
		out = append(out, k)
	}
	return out
}
