package route

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/FxIvan/bootstrap"
	"github.com/FxIvan/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
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

func deriveKey(email, salt string) ([]byte, error) {
	return scrypt.Key([]byte(email), []byte(salt), 32768, 8, 1, 32)
}
// Función para encriptar el correo electrónico
func encryptText(plainText string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    plaintext := []byte(plainText)
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]

    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptText(encryptedText string, key []byte) (string, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    if len(ciphertext) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return string(ciphertext), nil
}

