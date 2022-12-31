package firebase

import (
	"context"
	"idea-training-version-go/internals/models"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	errors "github.com/pkg/errors"
	"google.golang.org/api/option"
)

var (
	client *auth.Client
	ctx    context.Context
)

func init() {
	// Get an auth client from the firebase app.
	ctx = context.Background()
	opt := option.WithCredentialsFile("/Users/hirokiseino/go/src/idea-training-version-go/internals/utils/firebase/firebaseConfig/devAccount.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		panic("Error initializing firebase app")
	}
	client, err = app.Auth(ctx)
	if err != nil {
		errors.Wrap(err, "Error getting firebase auth client")
	}
}

func CreateUserInFirebase(email, password, firstName, lastName string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(firstName + " " + lastName).
		Disabled(false)

	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return "", errors.Wrap(err, "Error creating user in firebase")
	}
	return u.UID, nil
}

func GetUserInFirebase(email string) (*auth.UserRecord, error) {
	u, err := client.GetUserByEmail(ctx, email)
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
	u, err := client.UpdateUser(ctx, existingUser.UID, params)
	if err != nil {
		return nil, errors.Wrap(err, "Error updating user in firebase")
	}
	return u, nil
}

func DeleteUserInFirebase(uid string) error {
	err := client.DeleteUser(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "Error deleting user in firebase")
	}
	return nil
}