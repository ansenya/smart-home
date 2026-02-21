package repositories

import (
	"context"
	"time"

	"llm-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	Get(ctx context.Context, sessionID uuid.UUID) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, sessionID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Table("panel_sessions").Create(session).Error
}

func (r *sessionRepository) Get(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	var session models.Session

	err := r.db.WithContext(ctx).
		Table("panel_sessions").
		Where("id = ? AND expires_at > ?", sessionID, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Table("panel_sessions").Save(session).Error
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Table("panel_sessions").
		Where("id = ?", sessionID).
		Delete(&models.Session{}).Error
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Table("panel_sessions").
		Where("user_id = ?", userID).
		Delete(&models.Session{}).Error
}
