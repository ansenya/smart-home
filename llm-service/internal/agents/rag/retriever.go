package rag

import "context"

type Document struct {
	ID      string
	Content string
	Score   float64
}

type Retriever interface {
	Retrieve(ctx context.Context, userID string, query string, topK int) ([]Document, error)
}
