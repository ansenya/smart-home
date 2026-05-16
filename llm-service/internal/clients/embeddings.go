package clients

import "context"

type Embedder interface {
	Embed(ctx context.Context, model string, input []string) ([][]float32, error)
}
