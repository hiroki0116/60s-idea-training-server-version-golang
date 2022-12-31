package main

import (
	"context"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/db"
	"idea-training-version-go/internals/routes"
	"idea-training-version-go/internals/services"
	"log"
	"os"

	errors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	server         *gin.Engine
	usercollection *mongo.Collection
	usercontroller controllers.IUserController
	userservice    services.IUserService
	userroute      routes.UserRoutes
	ctx            context.Context
	err            error
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
	usercollection = db.MongoDB.Database("60s-idea-training").Collection("users")
	usercontroller = controllers.NewUserController(usercollection, ctx)
	userservice = services.NewUserService(usercontroller)
	userroute = *routes.NewUserRoutes(userservice)
	server = gin.Default()
}

func main() {
	defer db.MongoDB.Disconnect(ctx)
	basepath := server.Group("/api")
	userroute.UserRoutes(basepath)
	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
