package middleware

import (
	"llm-service/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SessionAuth(repo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("sid")

		if err != nil || sessionID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sid required"})
			return
		}

		uid, err := uuid.Parse(sessionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid sid"})
			return
		}

		session, err := repo.Get(c.Request.Context(), uid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session not found or expired"})
			return
		}

		c.Set("user_id", session.UserID)
		c.Next()
	}
}
