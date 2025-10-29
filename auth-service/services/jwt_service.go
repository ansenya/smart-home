package services

import (
	"auth-server/models"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"math/big"
	"os"
	"time"
)

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	EmailConfirmed bool   `json:"email_confirmed"`
	jwt.RegisteredClaims
}

type Jwks struct {
	Keys []Jwk `json:"keys"`
}

type Jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwtService struct {
	accessPrivateKey, refreshPrivateKey *rsa.PrivateKey
	accessPublicKey, refreshPublicKey   *rsa.PublicKey
	refreshKID                          string
	accessTokenDuration                 time.Duration
	refreshTokenDuration                time.Duration
}

func (s *jwtService) GenerateAccessToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:         user.ID,
		Email:          user.Email,
		EmailConfirmed: user.Confirmed,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-server",
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.accessPrivateKey)
}

func (s *jwtService) GenerateRefreshToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-server",
			ID:        uuid.NewString(),
			Subject:   user.ID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.refreshKID
	return token.SignedString(s.refreshPrivateKey)
}

func (s *jwtService) ValidateAccessToken(token string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessPublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(token string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshPublicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

func (s *jwtService) GetAccessTokenDuration() time.Duration {
	return s.accessTokenDuration
}

func (s *jwtService) GetRefreshTokenDuration() time.Duration {
	return s.refreshTokenDuration
}

func (s *jwtService) GenerateJwks() Jwks {
	n := base64.RawURLEncoding.EncodeToString(s.refreshPublicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.refreshPublicKey.E)).Bytes())

	return Jwks{
		Keys: []Jwk{
			{
				Kty: "RSA",
				Use: "sig",
				Kid: s.refreshKID,
				Alg: "RS256",
				N:   n,
				E:   e,
			},
		},
	}
}

func loadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPrivateKeyFromPEM(keyData)
}

func loadPublicKey(filename string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

func computeKID(pubKey *rsa.PublicKey) string {
	nBytes := pubKey.N.Bytes()
	eBytes := big.NewInt(int64(pubKey.E)).Bytes()
	data := append(nBytes, eBytes...)

	hash := sha256.Sum256(data)
	kid := base64.RawURLEncoding.EncodeToString(hash[:])
	return kid
}

func NewJwtService() (JWTService, error) {
	accessPrivateKey, err := loadPrivateKey("keys/access_private.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot load private key: %w", err)
	}
	accessPublicKey, err := loadPublicKey("keys/access_public.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot load public key: %w", err)
	}

	refreshPrivateKey, err := loadPrivateKey("keys/refresh_private.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot load private key: %w", err)
	}
	refreshPublicKey, err := loadPublicKey("keys/refresh_public.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot load public key: %w", err)
	}

	return &jwtService{
		accessPrivateKey:     accessPrivateKey,
		accessPublicKey:      accessPublicKey,
		refreshPrivateKey:    refreshPrivateKey,
		refreshPublicKey:     refreshPublicKey,
		refreshKID:           computeKID(refreshPublicKey),
		accessTokenDuration:  15 * time.Minute,
		refreshTokenDuration: 30 * 24 * time.Hour,
	}, nil
}
