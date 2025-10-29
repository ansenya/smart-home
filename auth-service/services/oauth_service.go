package services

import (
	"auth-server/models"
	"auth-server/repository"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/url"
	"time"
)

type oauthService struct {
	oauthClientsRepository repository.OauthClientsRepository
	sessionRepository      repository.SessionRepository
	oauthCodeService       TemporaryCodeService
}

func (s *oauthService) Authorize(queries models.OAuthRequest, sid string) (string, error) {
	oauthClient, err := s.oauthClientsRepository.GetByID(queries.ClientID)
	if err != nil || oauthClient == nil {
		return "", err
	}
	if !oauthClient.Enabled {
		return "", fmt.Errorf("client not enabled")
	}

	session, err := s.sessionRepository.GetByID(sid)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", fmt.Errorf("session not found")
	}

	code := uuid.New().String()
	m, err := json.Marshal(models.OAuthData{
		ClientID: oauthClient.ID,
		UserID:   session.UserID,
		Code:     code,
	})
	if err != nil {
		return "", err
	}
	if err := s.oauthCodeService.Save(code, string(m), time.Minute*15); err != nil {
		return "", err
	}

	redirectURL := fmt.Sprintf(
		"%s?code=%s&client_id=%s&state=%s&scope=%s",
		oauthClient.RedirectURI,
		url.QueryEscape(code),
		url.QueryEscape(queries.ClientID),
		url.QueryEscape(queries.State),
		url.QueryEscape(queries.Scope),
	)
	return redirectURL, nil
}

func (s *oauthService) GetByID(id string) (*models.OauthClient, error) {
	return s.oauthClientsRepository.GetByID(id)
}

func (s *oauthService) GetByName(name string) (*models.OauthClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewOauthClientsService(db *gorm.DB, redis *redis.Client) OauthService {
	return &oauthService{
		oauthClientsRepository: repository.NewOauthClientsRepository(db),
		sessionRepository:      repository.NewSessionRepository(db),
		oauthCodeService:       NewOauthCodeService(redis),
	}
}
