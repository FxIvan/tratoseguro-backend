package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUser := os.Getenv("MONGO_USER")
	dbPass := os.Getenv("MONGO_PASS")
	dbHost := os.Getenv("MONGO_HOST")
	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s", dbUser, dbPass, dbHost)

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

	router := gin.Default()

	router.POST("/", postBooks)

	coll := client.Database("db").Collection("books")
	doc := Book{Title: "Atonement", Author: "Ian McEwan"}

	result, err := coll.InsertOne(context.TODO(), doc)

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	/*
		coll := client.Database("sample_mflix").Collection("movies")
		title := "Back to the Future"
		var result bson.M
		err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the title %s\n", title)
			return
		}
		if err != nil {
			panic(err)
		}
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
	*/
	router.Run(":8080")
}

type SimpleUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var library []Book

func postBooks(c *gin.Context) {

	var newBook Book

	if err := c.BindJSON((&newBook)); err != nil {
		return
	}

	fmt.Print(newBook)

}
