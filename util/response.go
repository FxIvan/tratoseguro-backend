package response

import (
	"github.com/gin-gonic/gin"
)

func ResponseStatus(status int, message string, c *gin.Context) {
	c.JSON(status, gin.H{"message": message, "code": status})
}