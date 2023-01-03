package controllers

import (
	"context"
	"idea-training-version-go/internals/models"

	errors "github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

type IUserController interface {
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(id primitive.ObjectID, user *models.User) error
	GetUserByID(id primitive.ObjectID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

func NewUserController(usercollection *mongo.Collection, ctx context.Context) IUserController {
	return &UserController{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (uc *UserController) CreateUser(user *models.User) (*models.User, error) {
	result, err := uc.usercollection.InsertOne(uc.ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "Error in InsertOne")
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to fetch user id")
	}
	user.ID = oid
	return user, err
}

func (uc *UserController) GetUserByID(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}
	if err := uc.usercollection.FindOne(uc.ctx, filter).Decode(&user); err != nil {
		return nil, errors.Wrap(err, "Error in FindOne")
	}
	return &user, nil
}

func (uc *UserController) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	filter := bson.D{
		bson.E{
			Key:   "email",
			Value: email,
		},
	}
	if err := uc.usercollection.FindOne(uc.ctx, filter).Decode(&user); err != nil {
		return nil, errors.Wrap(err, "Error in FindOne")
	}
	return &user, nil
}

func (uc *UserController) UpdateUser(id primitive.ObjectID, user *models.User) error {
	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: id,
		},
	}

	if result, _ := uc.usercollection.UpdateOne(uc.ctx, filter, bson.M{"$set": user}); result.MatchedCount != 1 {
		return errors.New("failed to update user. User not found")
	}
	return nil
}
