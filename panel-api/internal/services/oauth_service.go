package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"panel-api/internal/models"
	"strings"
	"time"
)

type OauthService interface {
	ExchangeCode(ctx context.Context, in *models.CodeExchange) (*models.Tokens, error)
	GetIdentity(ctx context.Context, accessToken string) (*models.User, error)
}

type oauthService struct {
	client *http.Client
	cfg    OauthConfig
}

type OauthConfig struct {
	BaseURL       string
	TokenEndpoint string
	UserEndpoint  string
	Timeout       time.Duration
}

func NewOauthService(cfg OauthConfig, client *http.Client) OauthService {
	if client == nil {
		client = &http.Client{
			Timeout: cfg.Timeout,
		}
	}

	return &oauthService{
		client: client,
		cfg:    cfg,
	}
}

func (s *oauthService) doJSON(
	ctx context.Context,
	req *http.Request,
	out any,
) error {
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("oauth error %d: %s", resp.StatusCode, body)
	}

	if out == nil {
		return nil
	}

	return json.Unmarshal(body, out)
}

func (s *oauthService) ExchangeCode(
	ctx context.Context,
	in *models.CodeExchange,
) (*models.Tokens, error) {
	if in.Code == "" {
		return nil, errors.New("code is required")
	}

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {in.Code},
		"client_id":     {in.ClientID},
		"client_secret": {"in.ClientSecret"}, // todo
		"redirect_uri":  {"in.RedirectURI"},  // todo
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.cfg.BaseURL+s.cfg.TokenEndpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var tokens models.Tokens
	if err := s.doJSON(ctx, req, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (s *oauthService) GetIdentity(
	ctx context.Context,
	accessToken string,
) (*models.User, error) {
	if accessToken == "" {
		return nil, errors.New("access token is required")
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		s.cfg.BaseURL+s.cfg.UserEndpoint,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	var user models.User
	if err := s.doJSON(ctx, req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
