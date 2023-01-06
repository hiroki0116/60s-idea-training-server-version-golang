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
	"time"

	"github.com/gin-contrib/cors"
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
	db.ConnectDB(os.Getenv("MONGO_URI"))
	// collections
	usercollection = db.MongoDB.Database("60s-idea-trainings").Collection("users")
	ideacollection = db.MongoDB.Database("60s-idea-trainings").Collection("idearecords")
	// controllers
	usercontroller = controllers.NewUserController(usercollection, ctx)
	ideacontroller = controllers.NewIdeaController(ideacollection, ctx)
	// services
	userservice = services.NewUserService(usercontroller)
	ideaservice = services.NewIdeaService(ideacontroller)
	// middleware
	requireauth = middleware.NewRequireAuth(usercontroller)
	// routes
	userroute = routes.NewUserRoutes(userservice, requireauth)
	idearoute = routes.NewIdeaRoutes(ideaservice, requireauth)

	server = gin.Default()
	// CORS
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://60s-idea-training-client-6vl01gmwz-hiroki0116.vercel.app", "https://60s-idea-training.vercel.app"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTION", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func main() {
	defer db.MongoDB.Disconnect(ctx)
	// setup routes
	basepath := server.Group("/api")
	userroute.UserRoutes(basepath)
	idearoute.IdeaRoutes(basepath)
	routes.UtilsRoutes(basepath)

	log.Fatalln(server.Run(":" + os.Getenv("PORT")))
}
