package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	MongoDB *mongo.Client
	err     error
)

func ConnectDB(MONGO_URI string) {
	mongoconn := options.Client().ApplyURI(string(os.Getenv("MONGO_URI")))
	ctx := context.Background()
	MongoDB, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = MongoDB.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Connected to MongoDB!")
}
