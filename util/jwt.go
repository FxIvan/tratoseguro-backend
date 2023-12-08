package util

import (
	"fmt"
	"log"

	"github.com/FxIvan/config"
	"github.com/golang-jwt/jwt/v5"
)
func VerifyToken(authToken string) bool {
	athorization , err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		    // Verificar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Método de firma inesperado: %v", token.Header["alg"])
				return nil, fmt.Errorf("Método de firma inesperado: %v", token.Header["alg"])
			}
			// Devolver la clave secreta
			return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return false
	}
	//Valida el tiempo de expiracion del token
	if !athorization.Valid {
		return false
	}

	return true
}