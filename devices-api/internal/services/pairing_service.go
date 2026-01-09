package services

import (
	"devices-api/internal/models"
	"devices-api/internal/repositories"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type PairingService interface {
	StartPairing(userID uuid.UUID) (code string, expires int, err error)
	ConfirmPairing(request *models.ConfirmPairingRequest) error
	GetStatus(code string) (string, error)
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

	ttl := 5 * time.Minute

	if err := p.cache.Set(code, userID, ttl); err != nil {
		return "", 0, err
	}

	return code, int(ttl.Seconds()), nil
}

func (p *pairingService) GetStatus(code string) (string, error) {
	_, err := p.cache.Get(code)
	if err != nil {
		return "", err
	}
	if _, err := p.cache.Get("done:" + code); err == nil {
		return "done", nil
	}
	return "waiting", nil
}

func (p *pairingService) ConfirmPairing(request *models.ConfirmPairingRequest) error {
	userID, err := p.cache.Get(request.Code)
	if err != nil {
		return errors.New("invalid or expired code")
	}

	if err := p.repo.RegisterDevice(userID, request.DeviceUID, request.MacAddress); err != nil {
		return err
	}

	_ = p.cache.Set("done:"+request.Code, uuid.Nil, time.Minute)
	_ = p.cache.Delete(request.Code)

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
