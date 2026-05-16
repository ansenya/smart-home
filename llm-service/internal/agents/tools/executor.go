package tools

import (
	"context"
	"fmt"
)

type Executor interface {
	Execute(ctx context.Context, name string, args []byte) (string, error)
}

type executor struct {
	reg Registry
}

func NewExecutor(reg Registry) Executor {
	return &executor{reg: reg}
}

func (e *executor) Execute(ctx context.Context, name string, args []byte) (string, error) {
	h, ok := e.reg.Get(name)
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return h(ctx, args)
}

// RuntimeAdapter bridges tools.Executor to runtime.ToolExecutor interface.
type RuntimeAdapter struct {
	Exec Executor
}

func (a *RuntimeAdapter) Call(ctx context.Context, name string, args string) (string, error) {
	return a.Exec.Execute(ctx, name, []byte(args))
}
