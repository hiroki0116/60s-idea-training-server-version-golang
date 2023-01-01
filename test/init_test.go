package test

import (
	"context"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/db"
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
	usercontroller controllers.IUserController
	userservice    services.IUserService
	userroute      routes.UserRoutes
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
	// controllers
	usercontroller = controllers.NewUserController(usercollection, ctx)
	// services
	userservice = services.NewUserService(usercontroller)
	// routes
	userroute = routes.NewUserRoutes(userservice)
	// server
	server = gin.Default()
}

func TestMain(m *testing.M) {
	defer db.MongoDB.Disconnect(ctx)
	// set endpoints
	basepath := server.Group("/api")
	userroute.UserRoutes(basepath)
	unitTest.SetRouter(server)

	// populate sample data
	// DO SOMETHING

	newLog := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)
	unitTest.SetLog(newLog)
	exitVal := m.Run()

	log.Println("Everything below run after ALL test")
	firebase.DeleteAllUsersInFirebase()
	DeleteSampleData(usercollection, ctx)
	os.Exit(exitVal)
}
