package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"auth-server/structs"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
	"time"
)

type authHandler struct {
	userService            services.UserService
	oauthClientsRepository services.OauthClientsService
	oauthCodesService      services.TemporaryCodeService
	jwtService             services.JWTService
}

func newAuthRouter(userService services.UserService, oauthClientsRepository services.OauthClientsService, oauthCodesService services.TemporaryCodeService, jwtService services.JWTService) *authHandler {
	return &authHandler{
		userService:            userService,
		oauthClientsRepository: oauthClientsRepository,
		oauthCodesService:      oauthCodesService,
		jwtService:             jwtService,
	}
}

func (h *authHandler) RegisterRoutes(c *gin.RouterGroup) {
	c.POST("/register", h.handleRegister)
	c.POST("/login", h.handleLogin)
	c.POST("/token", h.handleToken)
	c.POST("/refresh", h.handleRefresh)

	//c.POST("/reset-password")
	//c.POST("/reset-password/change-password")

	//codeRequestGroup := r.Group("/", svc.JWTAuthMiddleware())
	//codeRequestGroup.POST("/auth/register/code/request", svc.CodeRequestHandler)
	//codeRequestGroup.POST("/auth/register/code/confirm", svc.CodeConfirmationHandler)
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
	var queries struct {
		ClientID string `form:"client_id" binding:"required"`
		State    string `form:"state" binding:"required"`
		Scope    string `form:"scope"`
	}
	if err := c.ShouldBindQuery(&queries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "queries are missing",
		})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.GetByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	if err := h.userService.IsPasswordCorrect(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	oauthClient, err := h.oauthClientsRepository.GetByID(queries.ClientID)
	if err != nil || oauthClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid client",
		})
		return
	}
	if !oauthClient.Enabled {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "client is not enabled",
		})
		return
	}

	code := uuid.New().String()
	m, err := json.Marshal(structs.AuthData{
		ClientID: oauthClient.ID,
		UserID:   user.ID,
		Code:     code,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to marshal auth data",
		})
		return
	}

	if err := h.oauthCodesService.Save(code, string(m), time.Minute*15); err != nil {
		log.Printf("failed to save code: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to save code: %v", err),
		})
		return
	}

	redirectURL := fmt.Sprintf(
		"%s?code=%s&client_id=%s&state=%s&scope=%s",
		oauthClient.RedirectURI,
		url.QueryEscape(code),
		url.QueryEscape(queries.ClientID),
		url.QueryEscape(queries.State),
		url.QueryEscape(queries.Scope),
	)

	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectURL,
	})
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

	claims, err := h.jwtService.ValidateToken(req.RefreshToken)
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
