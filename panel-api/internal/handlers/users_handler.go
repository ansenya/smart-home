package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"panel-api/internal/config"
	"panel-api/internal/models"
	"panel-api/internal/services"

	"github.com/gin-gonic/gin"
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
	group.GET("/me")
	group.POST("/exchange-code", h.ExchangeCode)
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

	c.SetCookie("sid", session.ID.String(), int(tokens.ExpiresIn), "/", "smarthome.hipahopa.ru", true, true)
	c.JSON(http.StatusOK, session)
}
