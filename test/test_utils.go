package test

import (
	"context"
	"fmt"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils/firebase"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PopulateUserSampleData(mongo *mongo.Collection, ctx context.Context) error {

	type FBUser struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	// create 10 users data
	for i := 0; i < 10; i++ {
		fb := &FBUser{
			FirstName: fmt.Sprintf("first_name_%d", i),
			LastName:  fmt.Sprintf("last_name_%d", i),
			Email:     fmt.Sprintf("test_email%v@test.com", i),
			Password:  fmt.Sprintf("test_password%v", i),
		}

		firebaseUID, err := firebase.CreateUserInFirebase(fb.Email, fb.Password, fb.FirstName, fb.LastName)
		if err != nil {
			log.Fatal("Error creating user in test firebase: ", err)
			return err
		}

		user := models.User{
			FirebaseUID: firebaseUID,
			Email:       fb.Email,
			FirstName:   fb.FirstName,
			LastName:    fb.LastName,
		}

		if _, err := mongo.InsertOne(ctx, user); err != nil {
			log.Fatal("Error inserting new user: ", err)
			return err
		}
	}

	return nil
}

func DeleteSampleData(mongo *mongo.Collection, ctx context.Context) error {
	filter := bson.D{{}}
	_, err := mongo.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal("Error deleting sample data: ", err)
		return err
	}
	return nil
}
