package route

import (
	"fmt"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	group.POST("/signup", func(c *gin.Context) {
		err := controllers.SignUp(c, &db)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
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

func ProfileUser(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup){
	group.GET("/profile", func(c *gin.Context){
		fmt.Println("profile")
	})
}
