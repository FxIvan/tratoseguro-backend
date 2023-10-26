package route

import (
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

	 // Parámetros de encriptación (ajusta según tus necesidades)
		salt := []byte("salt")  // Sal para derivación de clave
		key, err := scrypt.Key([]byte(request.Email), salt, 32768, 8, 1, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Verificar si el correo electrónico ya existe en la base de datos
		encryptedEmail, err := encryptEmail(request.Email, key)
		encryptedPassword, err := encryptEmail(request.Password, key)

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



		// Guardar el nuevo usuario en la base de datos
		_ , err = db.Collection("users").InsertOne(cntx, userEncrypt)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "User created"})
	})
}

// Función para encriptar el correo electrónico
func encryptEmail(email string, key []byte) ([]byte, error) {
	// Utiliza la clave derivada para encriptar el correo electrónico
	encryptedEmail, err := scrypt.Key([]byte(email), key, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return encryptedEmail, nil
}

// Función para desencriptar el correo electrónico cuando sea necesario
func decryptEmail(encryptedEmail []byte, key []byte) (string, error) {
	// Utiliza la clave derivada para desencriptar el correo electrónico
	email, err := scrypt.Key(encryptedEmail, key, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return string(email), nil
}
