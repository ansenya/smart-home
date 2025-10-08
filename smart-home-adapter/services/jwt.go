package services

//
//import (
//	"crypto/rsa"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt/v5"
//	"net/http"
//	"os"
//	"strings"
//	"time"
//)
//
//type Claims struct {
//	UserID         string `json:"user_id"`
//	Email          string `json:"email"`
//	EmailConfirmed bool   `json:"email_confirmed"`
//	jwt.RegisteredClaims
//}
//
//func loadPublicKey(filename string) (*rsa.PublicKey, error) {
//	keyData, err := os.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//
//	return jwt.ParseRSAPublicKeyFromPEM(keyData)
//}
//
//func (s *DevicesService) ValidateToken(tokenString string) (*Claims, error) {
//	claims := &Claims{}
//	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
//			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//		}
//		return s.accessPublicKey, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	if !token.Valid {
//		return nil, fmt.Errorf("invalid token")
//	}
//
//	if claims.ExpiresAt.Before(time.Now()) {
//		return nil, fmt.Errorf("token is expired")
//	}
//
//	return claims, nil
//}
//
//func (s *DevicesService) JWTAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.Request.Header.Get("Authorization")
//		if authHeader == "" {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//				"error": "No Authorization header found",
//			})
//			return
//		}
//
//		parts := strings.SplitN(authHeader, " ", 2)
//		if !(len(parts) == 2 && parts[0] == "Bearer") {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//				"error": "Invalid Authorization header",
//			})
//			return
//		}
//
//		tokenString := parts[1]
//		claims, err := s.ValidateToken(tokenString)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//				"error": "Invalid token",
//			})
//			return
//		}
//
//		c.Set("userID", claims.UserID)
//		c.Set("userEmail", claims.Email)
//		c.Set("emailConfirmed", claims.EmailConfirmed)
//
//		c.Next()
//	}
//}
//
//func (s *DevicesService) EmailConfirmedAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		s.JWTAuthMiddleware()(c)
//
//		if c.IsAborted() {
//			return
//		}
//
//		emailConfirmed, exists := c.Get("emailConfirmed")
//		if !exists {
//			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
//				"error": "Email confirmation status not found",
//			})
//			return
//		}
//
//		confirmed, ok := emailConfirmed.(bool)
//		if !ok || !confirmed {
//			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
//				"error": "Email confirmation required for this endpoint",
//			})
//			return
//		}
//
//		c.Next()
//	}
//}
