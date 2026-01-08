package services

import (
	"devices-api/internal/repositories"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type PairingService interface {
	StartPairing(userID uuid.UUID) (code string, expires int, err error)
	ConfirmPairing(code string, deviceUID string) error
}

type pairingService struct {
	repo  repositories.PairingRepository
	cache repositories.PairingCache
}

func NewPairingService(repo repositories.PairingRepository, cache repositories.PairingCache) PairingService {
	return &pairingService{
		repo:  repo,
		cache: cache,
	}
}

func (p *pairingService) StartPairing(userID uuid.UUID) (string, int, error) {
	code := generateCode(6)

	ttl := 2 * time.Minute

	if err := p.cache.Set(code, userID, ttl); err != nil {
		return "", 0, err
	}

	return code, int(ttl.Seconds()), nil
}

func (p *pairingService) ConfirmPairing(code string, deviceUID string) error {
	userID, err := p.cache.Get(code)
	if err != nil {
		return errors.New("invalid or expired code")
	}

	deviceUUID, err := uuid.Parse(deviceUID)
	if err != nil {
		return errors.New("invalid device uid")
	}

	manufactured, err := p.repo.FindManufacturedByMAC(deviceUID)
	if err != nil || manufactured.Registered {
		return errors.New("device not available")
	}

	if err := p.repo.RegisterDevice(deviceUUID, userID); err != nil {
		return err
	}

	_ = p.cache.Delete(code)

	return nil
}

func generateCode(n int) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
