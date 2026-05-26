package handlers

import (
	"encoding/json"
	"llm-service/internal/dto"
	"llm-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type messagesHandler struct {
	chatService services.ChatService
}

func (h *messagesHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/messages", h.GetHistory)
	rg.POST("/messages", h.SendMessage)
}

func (h *messagesHandler) GetHistory(c *gin.Context) {
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	msgs, err := h.chatService.GetHistory(c, chatID, userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := dto.MessageListResponse{
		Messages: make([]dto.MessageResponse, 0, len(msgs)),
		HasMore:  len(msgs) == limit,
	}

	for _, msg := range msgs {
		response.Messages = append(response.Messages, dto.MessageResponse{
			ID:           msg.ID,
			Role:         string(msg.Role),
			Content:      msg.Content,
			ModelName:    msg.ModelName,
			InputTokens:  msg.InputTokens,
			OutputTokens: msg.OutputTokens,
			ToolCallID:   msg.ToolCallID,
			ToolName:     msg.ToolName,
			ToolArgs:     rawOrNil(msg.ToolArgs),
			ToolResult:   rawOrNil(msg.ToolResult),
			Status:       string(msg.Status),
			CreatedAt:    msg.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *messagesHandler) SendMessage(c *gin.Context) {
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

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	isStream := c.DefaultQuery("stream", "false") == "true"

	if isStream {
		h.sendMessageStream(c, chatID, userID, req.Content)
		return
	}

	resp, err := h.chatService.SendMessage(c.Request.Context(), chatID, userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	//c.JSON(http.StatusOK, dto.MessageResponse{
	//	ID:           msg.ID,
	//	Role:         string(msg.Role),
	//	Content:      msg.Content,
	//	ModelName:    msg.ModelName,
	//	InputTokens:  msg.InputTokens,
	//	OutputTokens: msg.OutputTokens,
	//	Status:       string(msg.Status),
	//	CreatedAt:    msg.CreatedAt,
	//})
	c.JSON(http.StatusOK, resp)
}

func (h *messagesHandler) sendMessageStream(c *gin.Context, chatID, userID uuid.UUID, content string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Transfer-Encoding", "chunked")

	tokenChan := make(chan string, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		if err := h.chatService.StreamResponse(c.Request.Context(), chatID, userID, content, tokenChan); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case token, ok := <-tokenChan:
			if !ok {
				c.SSEvent("done", dto.StreamChunk{Done: true})
				c.Writer.Flush()
				return
			}
			c.SSEvent("message", dto.StreamChunk{Token: token, Done: false})
			c.Writer.Flush()

		case err, ok := <-errChan:
			if !ok {
				// goroutine returned without error; keep draining tokenChan
				// until it closes. Nil-out errChan so this case stops firing.
				errChan = nil
				continue
			}
			c.SSEvent("error", dto.ErrorResponse{Error: err.Error()})
			c.Writer.Flush()
			return

		case <-c.Request.Context().Done():
			return
		}
	}
}

func NewMessagesHandler(svcs *services.Container) Handler {
	return &messagesHandler{
		chatService: svcs.ChatService,
	}
}

func rawOrNil(p *json.RawMessage) json.RawMessage {
	if p == nil {
		return nil
	}
	return *p
}
