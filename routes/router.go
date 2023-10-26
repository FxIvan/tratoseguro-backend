package route

import (
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	publicRoute := gin.Group("/api/v1")

	SignUp(env, timeout, db, publicRoute)
	//SignIn(env, timeout, db, publicRoute)

}
