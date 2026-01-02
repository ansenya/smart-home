package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"panel-api/internal/models"
	"strings"
	"time"
)

type OauthService struct {
	client *http.Client

	baseUrl string
}

func NewOauthService() *OauthService {
	return &OauthService{
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        5,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		},
		baseUrl: "https://id.smarthome.hipahopa.ru",
	}
}

func (s *OauthService) ExchangeCode(ctx context.Context, codeExchange *models.CodeExchange) (*models.Tokens, error) {
	values := url.Values{}
	values.Set("code", codeExchange.Code)
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", codeExchange.ClientID)

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		s.baseUrl,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		text, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed exchanging: %s", text)
	}

	var tokenResp models.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}
