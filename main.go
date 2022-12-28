package main

import (
	"context"
	"idea-training-version-go/internals/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	server *gin.Engine
	ctx    context.Context
)

func init() {
	ctx = context.TODO()
	if os.Getenv("STAGE") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	// collections
	db.ConnectDB(os.Getenv("MONGO_URI"))
	server = gin.Default()

}

func main() {
	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
