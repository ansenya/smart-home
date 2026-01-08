package repositories

import (
	"devices-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Get(sessionID uuid.UUID) (*models.Session, error)
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

func (r *sessionRepository) Get(sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session
	return &session, r.db.Where("id = ?", sessionID).First(&session).Error
}

func (r *sessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}
