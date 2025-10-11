package repository

import (
	"auth-server/models"
	"errors"
	"gorm.io/gorm"
)

type oauthClientsRepository struct {
	db *gorm.DB
}

func (r *oauthClientsRepository) GetByID(id string) (*models.OauthClient, error) {
	var client models.OauthClient
	if err := r.db.Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func NewOauthClientsRepository(db *gorm.DB) OauthClientsRepository {
	return &oauthClientsRepository{db: db}
}
