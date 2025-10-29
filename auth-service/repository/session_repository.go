package repository

import (
	"auth-server/models"
	"errors"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func (r *sessionRepository) Create(session *models.Session) error {
	if err := r.db.Create(session).Error; err != nil {
		return err
	}
	return nil
}

func (r *sessionRepository) GetByID(id string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("id = ?", id).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}
