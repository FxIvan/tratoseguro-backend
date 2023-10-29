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

	 	//Creamos la key a traves de un algoritmo de encriptacion, que la sacamos de la libreria scrypt
		 salt := "salt" // Cambia la "sal" a una cadena de texto
		 key, err := deriveKey(request.Email, salt)
		 if err != nil {
			 c.JSON(400, gin.H{"error": err.Error()})
			 return
		 } // salt es una variable aleatoria que se utiliza como "semilla" para generar la clave derivada
		//key, err := scrypt.Key([]byte(request.Email), salt, 32768, 8, 1, 32) // Clave derivada 32768 es la iteracion. 8 es el tamaño de la memoria. 1 es el paralelismo. 32 es el tamaño de la clave
		
		/*if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}*/

		// Encrypamos email y password
		encryptedEmail, err := encryptText(request.Email, key)
		encryptedPassword, err := encryptText(request.Password, key)

		userEncrypt := domain.SignupRequestEncrypt{
			Email:    encryptedEmail,
			Password: encryptedPassword,
		}

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"email": encryptedEmail}
		err = db.Collection("users").FindOne(cntx, filter).Decode(&userEncrypt)
		
		if err == nil {
			// El correo electrónico ya está en uso
			c.JSON(400, gin.H{"error": "Email already in use"})
			return
		}else if err != mongo.ErrNoDocuments {
			// Ocurrió un error diferente a "documento no encontrado"
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		result,err := decryptText(encryptedEmail, key)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		
		fmt.Println(result)


		// Guardar el nuevo usuario en la base de datos
		_ , err = db.Collection("users").InsertOne(cntx, userEncrypt)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

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

