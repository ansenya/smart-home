package repositories

import (
	"panel-api/internal/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) Get(sessionID string) (*models.Session, error) {
	var session models.Session
	return &session, r.db.Where("session_id = ?", sessionID).First(&session).Error
}
