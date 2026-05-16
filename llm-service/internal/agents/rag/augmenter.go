package rag

import "llm-service/internal/clients"

type Augmenter interface {
	Augment(messages []clients.Message, docs []Document) []clients.Message
}
