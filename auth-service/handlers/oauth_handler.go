package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type oauthHandler struct {
	userService            services.UserService
	oauthClientsRepository services.OauthService
	oauthCodesService      services.TemporaryCodeService
	jwtService             services.JWTService
}

func (h *oauthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/authorize", h.Authorize)
	rg.POST("/token", h.Token)
	rg.POST("/refresh", h.Refresh)
	rg.GET("/jwks", h.JWKs)
}

func newOAuthHandler(db *gorm.DB, redis *redis.Client) (*oauthHandler, error) {
	jwtService, err := services.NewJwtService()
	if err != nil {
		return nil, err
	}
	return &oauthHandler{
		userService:            services.NewUserService(db),
		oauthClientsRepository: services.NewOauthClientsService(db, redis),
		oauthCodesService:      services.NewOauthCodeService(redis),
		jwtService:             jwtService,
	}, nil
}

func (h *oauthHandler) Authorize(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var queries models.OAuthRequest
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

func (h *oauthHandler) Token(c *gin.Context) {
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

	var oauthData models.OAuthData
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

func (h *oauthHandler) Refresh(c *gin.Context) {
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

func (h *oauthHandler) JWKs(c *gin.Context) {
	jwks := h.jwtService.GenerateJwks()
	c.JSON(http.StatusOK, jwks)
}
