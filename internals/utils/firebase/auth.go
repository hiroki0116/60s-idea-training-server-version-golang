package firebase_auth

import (
	"context"
	"idea-training-version-go/internals/models"
	"os"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	errors "github.com/pkg/errors"
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
		errors.Wrap(err, "Error initializing firebase app")
	}
	if client, err = app.Auth(context.Background()); err != nil {
		errors.Wrap(err, "Error getting firebase auth client")
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
		return nil, errors.Wrap(err, "Error creating user in firebase")
	}
	return u, nil
}

func GetUserInFirebase(email string) (*auth.UserRecord, error) {
	u, err := client.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user in firebase")
	}
	return u, nil
}

func UpdateUserInFirebase(user *models.User) (*auth.UserRecord, error) {
	params := (&auth.UserToUpdate{}).
		Email(user.Email).
		DisplayName(user.FirstName + " " + user.LastName).
		Disabled(false)
	existingUser, err := GetUserInFirebase(user.Email)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user in firebase")
	}
	u, err := client.UpdateUser(context.Background(), existingUser.UID, params)
	if err != nil {
		return nil, errors.Wrap(err, "Error updating user in firebase")
	}
	return u, nil
}

func DeleteUserInFirebase(uid string) error {
	err := client.DeleteUser(context.Background(), uid)
	if err != nil {
		return errors.Wrap(err, "Error deleting user in firebase")
	}
	return nil
}
