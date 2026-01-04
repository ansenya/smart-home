package handlers

import (
	"net/http"
	"panel-api/internal/models"
	"panel-api/internal/services"

	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	oauthService *services.OauthService
}

func newUsersHandler() *usersHandler {
	return &usersHandler{
		oauthService: services.NewOauthService(),
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
