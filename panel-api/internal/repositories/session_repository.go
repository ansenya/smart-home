package repositories

import (
	"panel-api/internal/models"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(session *models.Session) error
	Get(sessionID string) (*models.Session, error)
	Update(session *models.Session) error
}
type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepository) Get(sessionID string) (*models.Session, error) {
	var session models.Session
	return &session, r.db.Where("id = ?", sessionID).First(&session).Error
}

func (r *sessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}
