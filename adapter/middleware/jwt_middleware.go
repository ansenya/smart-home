package middleware

import (
	"adapter/utils"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

type JSONWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	cachedJWKS   *JWKS
	cachedJWKSAt time.Time
	jwksTTL      = 5 * time.Minute
	jwksMutex    sync.Mutex
)

func fetchJWKS(url string) (*JWKS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func getCachedJWKS(url string) (*JWKS, error) {
	jwksMutex.Lock()
	defer jwksMutex.Unlock()

	if cachedJWKS != nil && time.Since(cachedJWKSAt) < jwksTTL {
		return cachedJWKS, nil
	}

	jwks, err := fetchJWKS(url)
	if err != nil {
		return nil, err
	}

	cachedJWKS = jwks
	cachedJWKSAt = time.Now()
	return jwks, nil
}

// --- Parse RSA public key from JWKS ---
func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	// безопасный способ преобразования e
	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	n := new(big.Int).SetBytes(nBytes)
	return &rsa.PublicKey{N: n, E: e}, nil
}

// --- Parse JWT using JWKS ---
func parseJWT(tokenString string, jwks *JWKS) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("token missing kid")
		}

		for _, key := range jwks.Keys {
			if key.Kid == kid {
				return parseRSAPublicKey(key.N, key.E)
			}
		}
		return nil, fmt.Errorf("unable to find key %s", kid)
	}

	return jwt.Parse(tokenString, keyFunc)
}

func JWTMiddleware() gin.HandlerFunc {
	jwksURL := utils.GetEnv("AUTH_SERVICE_BASE_URL", "https://api.id.smarthome.hipahopa.ru/oauth/jwks")
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}
		accessToken := parts[1]

		jwks, err := getCachedJWKS(jwksURL)
		if err != nil {
			log.Println("Failed to fetch JWKS:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to fetch JWKS"})
			return
		}

		token, err := parseJWT(accessToken, jwks)
		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sub claim"})
			return
		}

		c.Set("userID", sub)
		c.Next()
	}
}
