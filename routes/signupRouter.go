package route

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/config"
	"github.com/FxIvan/domain"
	response "github.com/FxIvan/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	group.POST("/signup", func(c *gin.Context) {
		var request domain.SignupRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var cntx = c.Request.Context()

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

		userEncrypt := domain.SignupRequestEncrypt{
			Email:    request.Email,
			Password: string(encryptedPassword),
		}

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"email": userEncrypt.Email}
		err = db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
		if err == nil {
			c.JSON(400, gin.H{"error": "Email already in use"})
			return
		}else if err != mongo.ErrNoDocuments {
			c.JSON(500, gin.H{"error --->": err.Error()})
			return
		}

		_ , err = db.Collection("users").InsertOne(cntx, userEncrypt)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("User created ->",userEncrypt.Email)
		c.JSON(200, gin.H{"message": "User created successfully"})
	})
}

func SignIn(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup){
	group.POST("/signin", func(c *gin.Context) {
		var request domain.SignInRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}


		var cntx = c.Request.Context()

		filter := bson.M{"email": request.Email}

		var userEncrypt domain.SigninRequestEncrypt

		err := db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
		if err != nil {
			response.ResponseStatus(400, "Unregistered user", c)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userEncrypt.Password), []byte(request.Password))

		if err != nil {
			response.ResponseStatus(400, "Email or password incorrect", c)
			return
		}

		expTimeMs, _ := strconv.Atoi(config.JWTExpirationMs)
		
		type JWTCustomClaims struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
			jwt.RegisteredClaims
		}
		
		claims := JWTCustomClaims{
			ID:    userEncrypt.ID,
			Email: userEncrypt.Email,
			Role:  "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expTimeMs) * time.Millisecond)),
			},
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if err != nil {
			log.Fatal(err)
		}
		
		fmt.Println("Token generated ->", token)
		
		jwtString, err := token.SignedString([]byte(config.JWTSecret))
		
		if err != nil {
			fmt.Println("Error generating token ->", err)
			response.ResponseStatus(400, "Error generating token", c)
			return
		}
		
		c.JSON(200, gin.H{"token": jwtString})
		
	})
}

