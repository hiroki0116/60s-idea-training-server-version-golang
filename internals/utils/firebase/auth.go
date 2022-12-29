package firebase_auth

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	opt    option.ClientOption
	config *firebase.Config
	client *auth.Client
)

func init() {
	// Initialize firebase admin
	if os.Getenv("STAGE") == "production" {
		opt = option.WithCredentialsFile("./firebaseConfig/prodAccount.json")
		config = &firebase.Config{ProjectID: "seconds-idea-training-prod"}
	} else {
		opt = option.WithCredentialsFile("./firebaseConfig/devAccount.json")
		config = &firebase.Config{ProjectID: "seconds-idea-training-dev"}
	}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("Error initializing firebase app: %v", err)
	}
	if client, err = app.Auth(context.Background()); err != nil {
		log.Fatalf("Error getting firebase auth client: %v", err)
	}
}

func CreateUserInFirebase(email, password, firstName, lastName string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(firstName + " " + lastName).
		Disabled(false)

	u, err := client.CreateUser(context.Background(), params)

	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
		return nil, err
	}
	return u, nil
}

func GetUserInFirebase(email string) (*auth.UserRecord, error) {
	u, err := client.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Fatalf("error getting firebase user: %v\n", err)
		return nil, err
	}
	return u, nil
}

func UpdateUserInFirebase(user *models.User) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		Email(user.Email).
		DisplayName(user.FirstName + " " + user.LastName).
		Disabled(false)
	existingUser, err := GetUserInFirebase(user.Email)
	u, err := client.UpdateUser(context.Background(), existingUser.UID, params)
	if err != nil {
		log.Fatalf("error updating firebase user: %v\n", err)
		return nil, err
	}
	return u, nil
}

func DeleteUserInFirebase(uid string) error {
	err := client.DeleteUser(context.Background(), uid)
	if err != nil {
		log.Fatalf("error deleting firebase user: %v\n", err)
		return err
	}
	return nil
}
