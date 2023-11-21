package domain

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}