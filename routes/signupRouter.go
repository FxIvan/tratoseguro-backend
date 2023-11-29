package route

import (
	"fmt"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/controllers"
	"github.com/FxIvan/domain"
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
	group.POST("/signin", func(c *gin.Context){
		Response := controllers.Login(c, &db)
		if(Response.Status == "error"){
			c.JSON(400, gin.H{"error": Response.Data})
			return
		}else{
			c.JSON(200, gin.H{"token": Response.Token})
			return
		}
	})
}

