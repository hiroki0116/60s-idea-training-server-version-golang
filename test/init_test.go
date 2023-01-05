package test

import (
	"context"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/db"
	"idea-training-version-go/internals/middleware"
	"idea-training-version-go/internals/routes"
	"idea-training-version-go/internals/services"
	"idea-training-version-go/internals/utils/firebase"
	"log"
	"os"
	"testing"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
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
)

func init() {
	ctx = context.TODO()

	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	// database connection
	db.ConnectDB("MONGO_URI")
	// collections
	usercollection = db.MongoDB.Database("60s-idea-training").Collection("users")
	ideacollection = db.MongoDB.Database("60s-idea-training").Collection("idearecords")
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
	// server
	server = gin.Default()
}

func TestMain(m *testing.M) {
	defer db.MongoDB.Disconnect(ctx)
	// set endpoints
	basepath := server.Group("/api")
	userroute.UserRoutes(basepath)
	idearoute.IdeaRoutes(basepath)
	unitTest.SetRouter(server)

	log.Println("\n==========================\nPopulating sample data first! Wait for a momment...\n==========================")
	firebase.DeleteAllUsersInFirebase()
	DeleteSampleData(usercollection, ctx)
	DeleteSampleData(ideacollection, ctx)
	PopulateUserSampleData(usercollection, ctx)
	PopulateIdeaSampleData(usercollection, ideacollection, ctx)

	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	unitTest.SetLog(newLog)
	exitVal := m.Run()

	log.Println("======================== Cleaning sample data! Wait for a momment... =======================")
	firebase.DeleteAllUsersInFirebase()
	DeleteSampleData(usercollection, ctx)
	DeleteSampleData(ideacollection, ctx)
	os.Exit(exitVal)
}
