package test

import (
	"context"
	"fmt"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils/firebase"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func PopulateIdeaSampleData(usercollection *mongo.Collection, ideacollection *mongo.Collection, ctx context.Context) error {
	// create one user sample data
	type FBUser struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	fb := &FBUser{
		FirstName: "first_name_100",
		LastName:  "last_name_100",
		Email:     "test_email100@test.com",
		Password:  "test_password100",
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
	result, err := usercollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal("Error inserting new user: ", err)
		return err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("Error fetching user id: ", err)
		return err
	}

	// create 10 ideas data
	for i := 0; i < 10; i++ {
		task := models.Idea{
			TopicTitle: fmt.Sprintf("test title %d", i),
			Ideas:      &[]string{fmt.Sprintf("idea_%d", i)},
			CreatedBy:  oid,
		}

		if _, err := ideacollection.InsertOne(ctx, task); err != nil {
			log.Fatal("Error inserting new task: ", err)
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

func GenerateJWTToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
