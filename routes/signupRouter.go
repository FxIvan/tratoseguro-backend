package route

import (
	"context"
	"fmt"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	group.POST("/signup", func(c *gin.Context) {
		
	})
}


func SignIn(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup){
	group.Get("/singin", func(c *gin.Context){
	})
}