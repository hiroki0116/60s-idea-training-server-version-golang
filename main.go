package main

import (
	"context"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/db"
	"idea-training-version-go/internals/middleware"
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
	ideacollection *mongo.Collection
	usercontroller controllers.IUserController
	ideacontroller controllers.IIdeaController
	userservice    services.IUserService
	ideaservice    services.IIdeaService
	requireauth    middleware.RequireAuth
	userroute      routes.UserRoutes
	idearoute      routes.IdeaRoutes
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
	ideacollection = db.MongoDB.Database("60s-idea-training").Collection("idearecords")
	usercontroller = controllers.NewUserController(usercollection, ctx)
	ideacontroller = controllers.NewIdeaController(ideacollection, ctx)
	userservice = services.NewUserService(usercontroller)
	ideaservice = services.NewIdeaService(ideacontroller)
	requireauth = middleware.NewRequireAuth(usercontroller)
	userroute = routes.NewUserRoutes(userservice, requireauth)
	idearoute = routes.NewIdeaRoutes(ideaservice, requireauth)
	server = gin.Default()
}

func main() {
	defer db.MongoDB.Disconnect(ctx)
	basepath := server.Group("/api")
	userroute.UserRoutes(basepath)
	idearoute.IdeaRoutes(basepath)
	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
