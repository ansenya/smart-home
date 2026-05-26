package handlers

import (
	"devices-api/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type streamHandler struct {
	service services.StreamService
}

func newStreamHandler(service services.StreamService) *streamHandler {
	return &streamHandler{service: service}
}

func (h *streamHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/stream", h.Stream)
}

func (h *streamHandler) Stream(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx := c.Request.Context()
	events, cancel := h.service.Subscribe(ctx, uid)
	defer cancel()

	// Initial event so the client can confirm the connection is alive.
	fmt.Fprintf(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			return false
		case ev, ok := <-events:
			if !ok {
				return false
			}
			data, err := json.Marshal(ev)
			if err != nil {
				return true
			}
			fmt.Fprintf(w, "event: %s\n", sseEventName(ev.Type))
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		}
	})

	_ = http.StatusOK
}

func sseEventName(t string) string {
	if t == "" {
		return "message"
	}
	return t
}
