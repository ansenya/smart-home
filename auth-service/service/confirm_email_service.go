package service

import (
	"auth-server/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func (s *AuthService) CodeRequestHandler(c *gin.Context) {
	jwtUserEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	user, err := s.repo.GetUserByEmail(fmt.Sprint(jwtUserEmail))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "cannot query database",
			})
			log.Printf("error while checking if user exists: %v", err)
			return
		}
	}
	if user.Confirmed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email is already confirmed",
		})
		return
	}

	err = s.sendEmailConfirmationCode(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send email confirmation code",
		})
		log.Printf("Failed to send email confirmation code: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "code sent",
	})
}

func (s *AuthService) CodeConfirmationHandler(c *gin.Context) {
	jwtUserEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	type CodeBody struct {
		Code string `json:"code"`
	}

	var body CodeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	user, err := s.repo.GetUserByEmail(fmt.Sprint(jwtUserEmail))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "cannot query database",
			})
			log.Printf("error while checking if user exists: %v", err)
			return
		}
	}
	if user.Confirmed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email is already confirmed",
		})
		return
	}

	correct, err := s.repo.CheckCodeWithEmail(body.Code, fmt.Sprint(jwtUserEmail))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "code is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot check code",
		})
		log.Printf("error while checking code: %v", err)
		return
	}

	if !correct {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "code is incorrect",
		})
		return
	}

	err = s.repo.SetUserStatusConfirmed(fmt.Sprint(jwtUserEmail))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to set user status confirmed",
		})
		log.Printf("Failed to set user status confirmed: %v\n", err)
		return
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		log.Printf("Failed to generate access token: %v\n", err)
		return
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		log.Printf("Failed to generate refresh token: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        user.Email,
		Confirmed:    user.Confirmed,
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	})
}
