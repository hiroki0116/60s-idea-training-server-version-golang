package main

import (
	"context"
	"idea-training-version-go/internals/db"
	"log"
	"os"

	errors "github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	server *gin.Engine
	ctx    context.Context
	err    error
)

func init() {
	ctx = context.TODO()
	if os.Getenv("STAGE") != "production" {
		if err = godotenv.Load(); err != nil {
			errors.Wrap(err, "Error loading .env file")
		}
	}
	// Connect to MongoDB
	db.ConnectDB(os.Getenv("MONGO_URI"))

	server = gin.Default()

}

func main() {
	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
