package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func YandexMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestId := c.GetHeader("X-Request-ID")
		if xRequestId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "missing X-Request-ID header",
			})
			return
		}

		c.Set("requestID", xRequestId)
		c.Set("userID", "ae6bbc75-7520-45bc-b27a-bac023407c59")

		c.Next()
	}
}
