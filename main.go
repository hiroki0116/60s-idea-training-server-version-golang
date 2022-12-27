package main

import (
	"context"
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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	server = gin.Default()
	server.GET("/api/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
}

func main() {
	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
