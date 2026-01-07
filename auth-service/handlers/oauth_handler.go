package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type oauthHandler struct {
	oauthService services.OauthService
	jwtService   services.JWTService
}

func (h *oauthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/authorize", h.Authorize)
	rg.POST("/token", h.Token)
	rg.POST("/refresh", h.Refresh)
	rg.GET("/userinfo", h.Userinfo)
	rg.POST("/endsession", h.EndSession)
	rg.GET("/jwks", h.JWKs)
}

func newOauthHandler(db *gorm.DB, redis *redis.Client) *oauthHandler {
	return &oauthHandler{
		oauthService: services.NewOauthClientsService(db, redis),
		jwtService:   services.NewJwtService(),
	}
}

func (h *oauthHandler) Authorize(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var params models.OauthRequest
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "required params missing",
		})
		return
	}

	redirectUri, err := h.oauthService.Authorize(params, sid)
	if err != nil {
		if errors.Is(err, services.ErrorInvalidSession) || errors.Is(err, services.ErrorClientNotEnabled) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		}
		log.Printf("failed to authorize: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authorize"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectUri,
	})
}

func (h *oauthHandler) Token(c *gin.Context) {
	var request models.AccessTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	tokenResponse, err := h.oauthService.GenerateTokens(request)
	if err != nil {
		if errors.Is(err, services.ErrorInvalidCode) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		log.Printf("failed to generate tokens: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (h *oauthHandler) Refresh(c *gin.Context) {
	var request models.RefreshTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	tokenResponse, err := h.oauthService.RefreshTokens(request)
	if err != nil {
		if errors.Is(err, services.ErrorInvalidCode) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		log.Printf("failed to generate tokens: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (h *oauthHandler) Userinfo(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authorization header required",
		})
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(auth, bearerPrefix) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authorization header",
		})
		return
	}

	accessToken := strings.TrimPrefix(auth, bearerPrefix)
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "access token required",
		})
		return
	}

	user, err := h.oauthService.GetUserinfo(accessToken)
	if err != nil {
		log.Printf("failed to get userinfo: %s", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid access token",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *oauthHandler) EndSession(c *gin.Context) {
	// todo
}

func (h *oauthHandler) JWKs(c *gin.Context) {
	jwks := h.jwtService.GenerateJwks()
	c.JSON(http.StatusOK, jwks)
}
