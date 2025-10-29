package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"auth-server/structs"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type authHandler struct {
	authService            services.AuthService
	userService            services.UserService
	oauthClientsRepository services.OauthService
	oauthCodesService      services.TemporaryCodeService
	jwtService             services.JWTService
}

func newAuthRouter(db *gorm.DB, userService services.UserService, oauthClientsRepository services.OauthService, oauthCodesService services.TemporaryCodeService, jwtService services.JWTService) (*authHandler, error) {
	authService, err := services.NewAuthService(db)
	if err != nil {
		return nil, err
	}
	return &authHandler{
		authService:            authService,
		userService:            userService,
		oauthClientsRepository: oauthClientsRepository,
		oauthCodesService:      oauthCodesService,
		jwtService:             jwtService,
	}, nil
}

func (h *authHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.handleRegister)
	rg.POST("/login", h.handleLogin)

	rg.POST("/authorize", h.Authorize)
	rg.POST("/token", h.handleToken)
	rg.POST("/refresh", h.handleRefresh)

	rg.GET("/jwks", h.handleJwks)
}

func (h *authHandler) Authorize(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var queries models.OauthRequest
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "queries are missing",
		})
		return
	}

	redirectUri, err := h.oauthClientsRepository.Authorize(queries, sid)
	if err != nil {
		log.Printf("failed to authorize: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectUri,
	})
}

func (h *authHandler) handleRegister(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	if !h.userService.IsPasswordValid(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	hash, err := h.userService.HashPassword(req.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: hash,
	}

	if err := h.userService.Create(&user); err != nil {
		log.Printf("error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *authHandler) handleLogin(c *gin.Context) {
	var request models.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.authService.Login(&request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email/password"})
		return
	}
	c.SetCookie("sid", session.ID, int(h.jwtService.GetRefreshTokenDuration().Milliseconds()), "/", "", false, false)
	c.JSON(http.StatusOK, session)
}

func (h *authHandler) handleToken(c *gin.Context) {
	var req struct {
		Code         string `form:"code" binding:"required"`
		ClientSecret string `form:"client_secret" binding:"required"`
		GrantType    string `form:"grant_type" binding:"required"`
		ClientID     string `form:"client_id" binding:"required"`
		RedirectURI  string `form:"redirect_uri" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	var oauthData structs.AuthData
	rawData, err := h.oauthCodesService.Get(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	if err := json.Unmarshal([]byte(rawData), &oauthData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get auth rawData",
		})
		return
	}

	if req.ClientID != oauthData.ClientID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.userService.GetByID(oauthData.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid request",
		})
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate access token",
		})
		return
	}
	refreshToken, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate refresh token",
		})
		return
	}

	if err := h.oauthCodesService.Delete(req.Code); err != nil {
		log.Printf("failed to delete code: %s", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type":    "Bearer",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    h.jwtService.GetAccessTokenDuration().Seconds(),
	})
}

func (h *authHandler) handleRefresh(c *gin.Context) {
	var req struct {
		GrantType    string `form:"grant_type" binding:"required"`
		RefreshToken string `form:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	claims, err := h.jwtService.ValidateAccessToken(req.RefreshToken)
	if err != nil {
		log.Printf("vv: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid request",
		})
		return
	}

	user, err := h.userService.GetByID(claims.UserID)
	if err != nil || user == nil {
		log.Printf("nf: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid request",
		})
		return
	}

	accessToken, err := h.jwtService.GenerateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	refreshTokenNew, err := h.jwtService.GenerateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token_type":    "Bearer",
		"access_token":  accessToken,
		"refresh_token": refreshTokenNew,
		"expires_in":    h.jwtService.GetAccessTokenDuration().Seconds(),
	})
}

func (h *authHandler) handleJwks(c *gin.Context) {
	jwks := h.jwtService.GenerateJwks()
	c.JSON(http.StatusOK, jwks)
}
