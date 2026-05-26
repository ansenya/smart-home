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

func (r *sessionRepository) Delete(id string) error {
	res := r.db.Where("id = ?", id).Delete(&models.Session{})
	return res.Error
}

func (r *sessionRepository) DeleteByUserID(userID string) error {
	res := r.db.Where("user_id = ?", userID).Delete(&models.Session{})
	return res.Error
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}
