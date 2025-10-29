package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"math/big"
	"net/http"
	"strings"
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

// Конвертация base64url в *rsa.PublicKey
func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	e := int(binary.BigEndian.Uint32(append(make([]byte, 4-len(eBytes)), eBytes...)))
	n := new(big.Int).SetBytes(nBytes)

	return &rsa.PublicKey{N: n, E: e}, nil
}

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
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		parts := strings.Split(authorization, " ")
		accessToken := parts[1]
		jwks, _ := fetchJWKS("https://api.smarthome.hipahopa.ru/auth/jwks")
		token, err := parseJWT(accessToken, jwks)
		if err != nil {
			log.Println(err)
		}
		userID, _ := token.Claims.GetSubject()

		c.Set("userID", userID)
		c.Next()
	}
}
