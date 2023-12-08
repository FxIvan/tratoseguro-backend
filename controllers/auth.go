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

type UserRegister struct {
	Email string `json:"email"`
	Password string `json:"password"`
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

func SignUp(c *gin.Context, db *mongo.Database) error{
		

	var request domain.SignupRequest

		// Validate input
		if err := c.ShouldBindJSON(&request); err != nil {
			 c.JSON(400, gin.H{"error": err.Error()})
			 return err
		}
		
		// Context for the database
		var cntx = c.Request.Context()

		// Encrypt password, pasa por el ecryptado como 10 veces
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return err
		}

		// Create user encrypt para despues guardar este en la base de datos
		userEncrypt := domain.SignupRequestEncrypt{
			Email:    request.Email,
			Password: string(encryptedPassword),
		}

		// Check if email already exists
		filter := bson.M{"email": userEncrypt.Email}
		err = db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
		if err == nil {
			c.JSON(400, gin.H{"error": "Email already in use"})
			return err
		}else if err != mongo.ErrNoDocuments {
			c.JSON(500, gin.H{"error --->": err.Error()})
			return err
		}

		// Insert user into database
		_ , err = db.Collection("users").InsertOne(cntx, userEncrypt)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return err
		}

		// Respuesta de usuario creado
		c.JSON(200, gin.H{"message": "User created successfully"})
		return nil
}