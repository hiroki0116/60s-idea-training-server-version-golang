package test

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteSampleData(mongo *mongo.Collection, ctx context.Context) error {
	filter := bson.D{{}}
	_, err := mongo.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal("Error deleting sample data: ", err)
		return err
	}
	return nil
}
