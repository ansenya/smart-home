package handlers

import "github.com/gin-gonic/gin"

type pairingHandler struct {
}

func newPairingHandler() *pairingHandler {
	return &pairingHandler{}
}

func (h *pairingHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/start", h.Start)
	group.GET("/confirm", h.Confirm)
}

func (h *pairingHandler) Start(c *gin.Context) {

}
func (h *pairingHandler) Confirm(c *gin.Context) {

}
