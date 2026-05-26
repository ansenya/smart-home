package services

import (
	"auth-server/repository"
	"auth-server/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ErrorResetTokenInvalid = errors.New("invalid or expired reset link")
	ErrorResetWeakPassword = errors.New("password too weak")
)

type PasswordResetService interface {
	Request(email, baseURL string) error
	Confirm(token, newPassword string) error
}

type passwordResetService struct {
	users           repository.UserRepository
	sessions        repository.SessionRepository
	codes           TemporaryCodeService
	passwordService PasswordService
	emailService    EmailService
}

func NewPasswordResetService(db *gorm.DB, redisClient *redis.Client) PasswordResetService {
	return &passwordResetService{
		users:           repository.NewUserRepository(db),
		sessions:        repository.NewSessionRepository(db),
		codes:           NewTemporaryCodeService(redisClient, "pwreset"),
		passwordService: NewPasswordService(),
		emailService:    NewEmailService(),
	}
}

// Request always returns nil to the caller so we don't leak which emails are
// registered. Failures are logged internally.
func (s *passwordResetService) Request(email, baseURL string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil
	}

	user, err := s.users.GetByEmail(email)
	if err != nil || user == nil {
		log.Printf("password reset: lookup failed or no user for %s (err=%v)", email, err)
		return nil
	}

	token, err := generateResetToken()
	if err != nil {
		log.Printf("password reset: token gen failed: %v", err)
		return nil
	}

	ttl := 30 * time.Minute
	if err := s.codes.Save(token, user.ID, ttl); err != nil {
		log.Printf("password reset: redis save failed: %v", err)
		return nil
	}

	if baseURL == "" {
		baseURL = utils.GetEnv("AUTH_WEB_URL", "https://id.smarthome.hipahopa.ru")
	}
	resetURL := fmt.Sprintf("%s/reset-password/confirm?token=%s", strings.TrimRight(baseURL, "/"), token)

	if err := s.emailService.SendPasswordReset(email, resetURL); err != nil {
		log.Printf("password reset: send failed: %v (token will still be usable: %s)", err, resetURL)
	}
	return nil
}

func (s *passwordResetService) Confirm(token, newPassword string) error {
	if token == "" {
		return ErrorResetTokenInvalid
	}
	if !s.passwordService.IsPasswordValid(newPassword) {
		return ErrorResetWeakPassword
	}

	userID, err := s.codes.Get(token)
	if err != nil || userID == "" {
		return ErrorResetTokenInvalid
	}

	hash, err := s.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.users.UpdatePassword(userID, hash); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorResetTokenInvalid
		}
		return fmt.Errorf("update password: %w", err)
	}

	if err := s.sessions.DeleteByUserID(userID); err != nil {
		log.Printf("password reset: failed to invalidate sessions: %v", err)
	}

	if err := s.codes.Delete(token); err != nil {
		log.Printf("password reset: failed to delete token: %v", err)
	}
	return nil
}

func generateResetToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
