package service

import (
	"auth-server/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/mail"
	"net/url"
)

func (s *AuthService) RegisterHandler(c *gin.Context) {
	var req models.RegistrationUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	// validate email
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email address",
		})
		return
	}
	// validate password
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters",
		})
		return
	}

	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists",
		})
		return
	}
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			log.Printf("error while checking if user exists: %v", err)
			return
		}
	}

	user, err := s.repo.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		log.Printf("Failed to create user: %v\n", err)
		return
	}

	err = s.sendEmailConfirmationCode(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send email confirmation code",
		})
		log.Printf("failed to send email confirmation code: %v\n", err)
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

	//c.JSON(http.StatusCreated, models.LoginResponse{
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//	TokenType:    "Bearer",
	//	Email:        "",
	//	Confirmed:    false,
	//	ExpiresIn:    0,
	//})
	c.SetCookie("access_token",
		accessToken,
		int(s.accessTokenDuration.Seconds()),
		"/",
		"",
		true,
		true)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(s.refreshTokenDuration.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"email":     user.Email,
		"confirmed": user.Confirmed,
	})
}

func (s *AuthService) LoginHandler(c *gin.Context) {
	var req models.LoginUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// this part is for yandex-smarthome auth (https://yandex.ru/dev/dialogs/smart-home/doc/ru/auth/how-it-works)
	if c.Query("smart-home") == "true" {
		code := uuid.New().String()
		clientId := c.Query("client_id")
		userId := user.ID
		state := c.Query("state")
		scope := c.Query("scope")

		err := s.repo.SetLoginDataByCode(code, clientId, userId, state, scope)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to set login data",
			})
			log.Printf("failed to set login data: %v\n", err)
			return
		}

		redirectURL := fmt.Sprintf(
			"https://social.yandex.net/broker/redirect?code=%s&client_id=%s&state=%s&scope=%s",
			url.QueryEscape(code),
			url.QueryEscape(clientId),
			url.QueryEscape(state),
			url.QueryEscape(scope),
		)

		c.JSON(http.StatusOK, gin.H{
			"redirect_url": redirectURL,
		})
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

	c.SetCookie("access_token",
		accessToken,
		int(s.accessTokenDuration.Seconds()),
		"/",
		"",
		true,
		true)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(s.refreshTokenDuration.Seconds()),
		"/",
		"",
		true,
		true,
	)

	//c.JSON(http.StatusOK, models.LoginResponse{
	//	TokenType:    "Bearer",
	//	AccessToken:  accessToken,
	//	RefreshToken: refreshToken,
	//	Email:        user.Email,
	//	Confirmed:    user.Confirmed,
	//	Name:         user.Name,
	//	ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	//})
	c.JSON(http.StatusOK, gin.H{
		"email":     user.Email,
		"confirmed": user.Confirmed,
	})
}

func (s *AuthService) RefreshTokenHandler(c *gin.Context) {
	if c.ContentType() != "application/x-www-form-urlencoded" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "unsupported_media_type"})
		return
	}

	grantType := c.PostForm("grant_type")
	if grantType != "refresh_token" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
		return
	}

	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "error_description": "Missing refresh_token"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.refreshPublicKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "invalid_grant",
			"error_description": "Invalid or expired refresh token",
		})
		return
	}

	user, err := s.repo.GetUserById(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "invalid_grant",
			"error_description": "User not found",
		})
		return
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	refreshTokenNew, err := s.generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	yandexQuery := c.Query("yandex-smart-home")
	isYandexRequest := yandexQuery != ""

	if isYandexRequest {
		c.JSON(http.StatusOK, models.LoginResponse{
			TokenType:    "Bearer",
			AccessToken:  accessToken,
			RefreshToken: refreshTokenNew,
			ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
		})
	} else {
		c.SetCookie(
			"access_token",
			accessToken,
			int(s.accessTokenDuration.Seconds()),
			"/",
			"",
			true,
			true,
		)

		c.SetCookie(
			"refresh_token",
			refreshTokenNew,
			int(s.refreshTokenDuration.Seconds()),
			"/",
			"",
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"email":     user.Email,
			"confirmed": user.Confirmed,
		})
	}
}

// GetTokenHandler is for yandex-smarthome auth (https://yandex.ru/dev/dialogs/smart-home/doc/ru/auth/how-it-works)
func (s *AuthService) GetTokenHandler(c *gin.Context) {
	var req models.TokenRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	loginData, err := s.repo.GetLoginDataByCode(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve login data",
		})
		log.Printf("Failed to retrieve login data: %v\n", err)
		return
	}

	if loginData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Login data not found",
		})
		return
	}

	if loginData.ClientID != req.ClientID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Login data not found",
		})
		return
	}

	user, err := s.repo.GetUserById(loginData.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
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

	if s.repo.DeleteLoginDataByCode(req.Code) != nil {
		log.Printf("Failed to delete login data: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type":    "Bearer",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    s.accessTokenDuration.Seconds(),
	})
}
