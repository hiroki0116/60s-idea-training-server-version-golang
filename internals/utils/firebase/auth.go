package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/joho/godotenv"
	errors "github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	Client *auth.Client
	ctx    context.Context
)

func init() {
	// Get an auth Client from the firebase app.
	ctx = context.Background()
	if os.Getenv("STAGE") == "test" {
		if err := godotenv.Load(".env.test"); err != nil {
			errors.Wrap(err, "Error loading .env.test file")
		}
	} else if os.Getenv("STAGE") == "production" {
	} else {
		if err := godotenv.Load(); err != nil {
			errors.Wrap(err, "Error loading .env file")
		}
	}
	opt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CRED")))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln("Error initializing firebase app", err)
	}
	Client, err = app.Auth(ctx)
	if err != nil {
		errors.Wrap(err, "Error getting firebase auth Client")
	}
}

func CreateUserInFirebase(email, password, firstName, lastName string) (string, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(firstName + " " + lastName).
		Disabled(false)

	u, err := Client.CreateUser(ctx, params)
	if err != nil {
		return "", errors.Wrap(err, "Error creating user in firebase")
	}
	return u.UID, nil
}

func GetUserInFirebase(email string) (*auth.UserRecord, error) {
	u, err := Client.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting user in firebase")
	}
	return u, nil
}

func DeleteUserInFirebase(uid string) error {
	err := Client.DeleteUser(ctx, uid)
	if err != nil {
		return errors.Wrap(err, "Error deleting user in firebase")
	}
	return nil
}

func DeleteAllUsersInFirebase() error {
	iter := Client.Users(ctx, "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("error listing fireabase users")
			return err
		}
		Client.DeleteUser(ctx, user.UID)
	}
	return nil
}
