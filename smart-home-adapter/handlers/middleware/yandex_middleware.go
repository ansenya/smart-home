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
		c.Set("userID", "b84ae96d-b311-49dd-b762-2ecec42c1bd4")

		c.Next()
	}
}
