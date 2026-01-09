package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"panel-api/internal/config"
	"panel-api/internal/models"
	"panel-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type usersHandler struct {
	oauthService services.OauthService
	usersService services.UsersService

	log *slog.Logger
}

func newUsersHandler(cfg *config.Container) *usersHandler {
	return &usersHandler{
		oauthService: cfg.Services.OauthService,
		usersService: cfg.Services.UsersService,
		log:          cfg.Log,
	}
}

func (h *usersHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/me", h.Me)
	group.POST("/logout", h.Logout)
	group.POST("/exchange-code", h.ExchangeCode)
}

func (h *usersHandler) Me(c *gin.Context) {
	sid, err := c.Cookie(SessionIDName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session cookie not found"})
		return
	}

	user, err := h.usersService.GetUserBySessionID(sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Info("user not found", sid)
			c.Status(http.StatusUnauthorized)
			return
		}
		h.log.Error("could not get user by session id", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *usersHandler) Logout(c *gin.Context) {
	sid, err := c.Cookie(SessionIDName)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := h.usersService.ExpireSession(sid); err != nil {
		h.log.Error("could not expire session", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	// todo: expire id.smarthome session

	c.SetCookie(SessionIDName, "", 0, "/", DomainName, true, true)
	c.Status(http.StatusOK)
}

func (h *usersHandler) ExchangeCode(c *gin.Context) {
	var request models.CodeExchange
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	tokens, err := h.oauthService.ExchangeCode(c, &request)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to exchange code: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user, err := h.oauthService.GetIdentity(c, tokens.AccessToken)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to get identity: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	session, err := h.usersService.CreateSession(user, tokens)
	if err != nil {
		h.log.Error(fmt.Sprintf("failed to create session: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.SetCookie(SessionIDName, session.ID.String(), int(tokens.ExpiresIn), "/", DomainName, true, true)
	c.JSON(http.StatusOK, session)
}
