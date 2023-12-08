package main

import (
	"context"
	"fmt"
	"log"
	"time"

	bootstrap "github.com/FxIvan/bootstrap"
	route "github.com/FxIvan/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	appEnv := bootstrap.NewEnv()
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s", appEnv.DbUser, appEnv.DbPass, appEnv.DbHost)
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	type Book struct {
		Title  string
		Author string
	}

	timeout := time.Duration(5) * time.Second

	gin := gin.Default()
	route.Setup(appEnv, timeout, *client.Database("books"), gin)
	gin.Run(":8080")
}

type SimpleUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
