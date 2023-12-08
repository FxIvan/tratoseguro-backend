package middleware

import (
	"fmt"
	"strings"

	"github.com/FxIvan/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

		fmt.Println("authToken", authToken)

		//Verificacion de Token
		authorize, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		if !authorize.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		
		fmt.Println("authorize", authorize)

	}
}