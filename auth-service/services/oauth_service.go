package services

import (
	"auth-server/models"
	"auth-server/repository"
	"fmt"
)

type oauthClientsService struct {
	oauthClientsRepository repository.OauthClientsRepository
}

func (s *oauthClientsService) GetByID(id string) (*models.OauthClient, error) {
	return s.oauthClientsRepository.GetByID(id)
}

func (s *oauthClientsService) GetByName(name string) (*models.OauthClient, error) {
	return nil, fmt.Errorf("not implemented")
}

func NewOauthClientsService(oauthClientsRepository repository.OauthClientsRepository) OauthClientsService {
	return &oauthClientsService{oauthClientsRepository: oauthClientsRepository}
}
