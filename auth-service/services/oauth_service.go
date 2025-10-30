package services

import (
	"auth-server/models"
	"auth-server/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"net/url"
	"time"
)

var (
	ErrorClientNotEnabled = errors.New("client not enabled")
	ErrorInvalidSession   = errors.New("invalid session")
	ErrorInvalidToken     = errors.New("invalid token")
	ErrorInvalidCode      = errors.New("invalid code")
)

type oauthService struct {
	userRepository         repository.UserRepository
	oauthClientsRepository repository.OauthClientsRepository
	sessionRepository      repository.SessionRepository
	jwtService             JWTService
	oauthCodeService       TemporaryCodeService
}

func (s *oauthService) Authorize(queries models.OauthRequest, sid string) (string, error) {
	oauthClient, err := s.oauthClientsRepository.GetByID(queries.ClientID)
	if err != nil || oauthClient == nil {
		return "", err
	}
	if !oauthClient.Enabled {
		return "", ErrorClientNotEnabled
	}

	session, err := s.sessionRepository.GetByID(sid)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", ErrorInvalidSession
	}

	code := uuid.New().String()
	m, err := json.Marshal(models.OauthData{
		ClientID: oauthClient.ID,
		UserID:   session.UserID,
		Code:     code,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal oauth data: %s", err)
	}
	if err := s.oauthCodeService.Save(code, string(m), time.Minute*15); err != nil {
		return "", fmt.Errorf("failed to save oauth data")
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

func (s *oauthService) GenerateTokens(request models.AccessTokenRequest) (*models.TokenResponse, error) {
	raw, err := s.oauthCodeService.Get(request.Code)
	if err != nil {
		return nil, ErrorInvalidCode
	}
	var oauthData models.OauthData
	if err := json.Unmarshal([]byte(raw), &oauthData); err != nil {
		return nil, fmt.Errorf("failed to unsmarshal oauth data: %s", err)
	}

	if request.ClientID != oauthData.ClientID {
		return nil, ErrorInvalidCode
	}

	user, err := s.userRepository.GetByID(oauthData.UserID)
	if err != nil {
		log.Printf("user not found by oauth data: %s\ndata: %s", err, raw)
		return nil, ErrorInvalidCode
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %s", err)
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %s", err)
	}

	if err := s.oauthCodeService.Delete(request.Code); err != nil {
		log.Printf("failed to delete code: %s", err)
	}

	return &models.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtService.GetAccessTokenDuration().Seconds(),
	}, nil
}

func (s *oauthService) RefreshTokens(request models.RefreshTokenRequest) (*models.TokenResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		log.Printf("failed to valudate refresh token: %s", err)
		return nil, ErrorInvalidToken
	}
	user, err := s.userRepository.GetByID(claims.Subject)
	if err != nil || user == nil {
		log.Printf("user not found by jwt subject id: %s\nsubject: %s", err, claims.Subject)
		return nil, ErrorInvalidToken
	}
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %s", err)
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %s", err)
	}
	return &models.TokenResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtService.GetAccessTokenDuration().Seconds(),
	}, nil
}

func NewOauthClientsService(db *gorm.DB, redis *redis.Client) OauthService {
	return &oauthService{
		userRepository:         repository.NewUserRepository(db),
		oauthClientsRepository: repository.NewOauthClientsRepository(db),
		sessionRepository:      repository.NewSessionRepository(db),
		jwtService:             NewJwtService(),
		oauthCodeService:       NewTemporaryCodeService(redis, "oauth"),
	}
}
