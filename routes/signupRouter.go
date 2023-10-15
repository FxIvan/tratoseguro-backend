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
		var request domain.SignupRequest

		err := c.ShouldBind(&request)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
		plain_password := []byte(request.Password)
		plain_email := []byte(request.Email)


		block, err := aes.NewCipher(key)

		if err != nil {
			panic(err.Error())
		}

		nonce := make([]byte, 12)


		if _,err := io.ReadFull(rand.Reader, nonce); err != nil {
			fmt.Println("ERROR ->",err.Error())
			panic(err.Error())
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}

		ciphertext := aesgcm.Seal(nil, nonce, plain_password, nil)
		ciphertext2 := aesgcm.Seal(nil, nonce, plain_email, nil)

		fmt.Printf("%x\n", ciphertext)
		fmt.Printf("%x\n", ciphertext2)

		user := domain.SignupRequestEncrypt{
			Email:    ciphertext,
			Password: ciphertext2,
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
