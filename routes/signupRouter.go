package route

import (
	"fmt"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/domain"
	response "github.com/FxIvan/util"
	"github.com/gin-gonic/gin"
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
		var resuest domain.SignInRequest

		if err := c.ShouldBindJSON(&resuest); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}


		var cntx = c.Request.Context()

		filter := bson.M{"email": resuest.Email}

		var userEncrypt domain.SigninRequestEncrypt

		err := db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
		if err != nil {
			response.ResponseStatus(400, "Unregistered user", c)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userEncrypt.Password), []byte(resuest.Password))

		if err != nil {
			response.ResponseStatus(400, "Email or password incorrect", c)
			return
		}

		fmt.Println("User logged ->",userEncrypt.Email)

	})
}

