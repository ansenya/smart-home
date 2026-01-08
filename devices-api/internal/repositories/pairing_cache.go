package repositories

import (
	"context"
	"devices-api/internal/infra/rds"
	"time"

	"github.com/google/uuid"
)

type PairingCache interface {
	Set(code string, userID uuid.UUID, ttl time.Duration) error
	Get(code string) (uuid.UUID, error)
	Delete(code string) error
}

type pairingService struct {
	cache rds.NamespacedRedis
}

func NewPairingCache(rdb rds.NamespacedRedis) PairingCache {
	return &pairingService{
		cache: rdb,
	}
}

func (p pairingService) Set(code string, userID uuid.UUID, ttl time.Duration) error {
	return p.cache.Set(context.Background(), code, userID.String(), ttl)
}

func (p pairingService) Get(code string) (uuid.UUID, error) {
	get, err := p.cache.Get(context.Background(), code)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(get)
}

func (p pairingService) Delete(code string) error {
	return p.cache.Del(context.Background(), code)
}
