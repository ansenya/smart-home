package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func getValueFromContext(c *gin.Context, key string) string {
	value, _ := c.Get(key)
	return fmt.Sprint(value)
}
