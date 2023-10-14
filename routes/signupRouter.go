package route

import (
	"context"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	group.POST("/signup", func(c *gin.Context) {
		var request domain.SignupRequest

		err := c.ShouldBind(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		user := domain.SignupRequest{
			Email:    request.Email,
			Password: request.Password,
		}

		collection := db.Collection("users")

		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User created successfully"})

	})

}
