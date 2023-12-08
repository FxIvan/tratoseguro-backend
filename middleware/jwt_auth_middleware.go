package middleware

import (
	"fmt"
	"strings"

	"github.com/FxIvan/util"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) != 2 || t[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		
		authToken := t[1]
		authorize := util.VerifyToken(authToken)

		if !authorize {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		
		fmt.Println("authorize", authorize)

	}
}