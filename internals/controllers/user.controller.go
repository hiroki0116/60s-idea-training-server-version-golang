package controllers

import (
	"context"
	"idea-training-version-go/internals/models"

	errors "github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

type IUserController interface {
	CreateUser(user *models.User) (*models.User, error)
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
