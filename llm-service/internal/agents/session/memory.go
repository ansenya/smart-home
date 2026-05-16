package session

// prompt-injection

type MemoryStore interface {
	GetFacts(userID string) ([]string, error)
	AddFact(userID string, fact string) error
	RemoveFact(userID string, fact string) error
}
