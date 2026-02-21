package handlers

import (
	"llm-service/internal/dto"
	"llm-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chatHandler struct {
	chatService services.ChatService
}

func (h *chatHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.CreateChat)
	rg.GET("", h.GetChats)

	rg.GET("/:id", h.GetByID)
	rg.PUT("/:id", h.UpdateChat)
	rg.DELETE("/:id", h.DeleteChat)
}

func (h *chatHandler) CreateChat(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var req dto.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	chat, err := h.chatService.CreateChat(c, userID, req.Model, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ChatResponse{
		ID:        chat.ID,
		Model:     chat.Model,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	})
}

func (h *chatHandler) GetChats(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	chats, err := h.chatService.GetUserChats(c, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.ChatListResponse{
		Chats: make([]dto.ChatResponse, 0, len(chats)),
	}

	for _, chat := range chats {
		response.Chats = append(response.Chats, dto.ChatResponse{
			ID:        chat.ID,
			Model:     chat.Model,
			Title:     chat.Title,
			CreatedAt: chat.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *chatHandler) GetByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid chat id"})
		return
	}

	chat, err := h.chatService.GetChat(c, chatID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "chat not found"})
		return
	}

	c.JSON(http.StatusOK, dto.ChatResponse{
		ID:        chat.ID,
		Model:     chat.Model,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	})
}

func (h *chatHandler) UpdateChat(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid chat id"})
		return
	}

	var req dto.UpdateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	chat, err := h.chatService.UpdateChat(c, chatID, userID, req.Title, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ChatResponse{
		ID:        chat.ID,
		Model:     chat.Model,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	})
}

func (h *chatHandler) DeleteChat(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid chat id"})
		return
	}

	if err := h.chatService.DeleteChat(c.Request.Context(), chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func getUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, nil
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, nil
	}

	return id, nil
}

func NewChatHandler(svcs *services.Container) Handler {
	return &chatHandler{
		chatService: svcs.ChatService,
	}
}
