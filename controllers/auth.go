package controllers

import (
	"strconv"
	"time"

	"github.com/FxIvan/config"
	"github.com/FxIvan/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Token string `json:"token"`
	Status string `json:"status"`
	Data string `json:"data"`
}

func Login(c *gin.Context, db *mongo.Database) Response{
	var request domain.SignInRequest

	var response string;

	if err := c.ShouldBindJSON(&request); err != nil {
		response = "Error en el request"
		salida := Response{
			Token: "",
			Status: "error",
			Data: response,
		}
		return salida
	}

	var cntx = c.Request.Context()

	filter := bson.M{"email": request.Email}

	var userEncrypt domain.SigninRequestEncrypt

	err := db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
	if err != nil{
		c.JSON(400, gin.H{"error": err.Error()})
		response = "La contraseña o el email son incorrectos"
		salida := Response{
			Token: "",
			Status: "error",
			Data: response,
		}
		return  salida
	}

	err = bcrypt.CompareHashAndPassword([]byte(userEncrypt.Password), []byte(request.Password))

	if err != nil {
		response = "La contraseña o el email son incorrectos"
		salida := Response{
			Token: "",
			Status: "error",
			Data: response,
		}
		return salida
	}

	expireTimeMs, _ :=  strconv.Atoi(config.JWTExpirationMs)

	type JWTCustomClaims struct {
		ID string `json:"id"`
		Email string `json:"email"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	claims := JWTCustomClaims{
		ID: userEncrypt.ID,
		Email: userEncrypt.Email,
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTimeMs) * time.Millisecond)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtString, err := token.SignedString([]byte(config.JWTSecret))

	if err != nil {
		response = "Error generando el token"
		salida := Response{
			Token: "",
			Status: "error",
			Data: response,
		}
		return salida
	}

	response = jwtString
	return Response{
		Token: response,
		Status: "success",
		Data: response,
	}

}