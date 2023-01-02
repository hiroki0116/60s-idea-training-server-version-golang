package controllers

import (
	"context"
	"errors"
	"idea-training-version-go/internals/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdeaController struct {
	ideacollection *mongo.Collection
	ctx            context.Context
}

type IIdeaController interface {
	CreateIdea(idea *models.Idea) (*models.Idea, error)
}

func NewIdeaController(ideacollection *mongo.Collection, ctx context.Context) IIdeaController {
	return &IdeaController{
		ideacollection: ideacollection,
		ctx:            ctx,
	}
}

func (ic *IdeaController) CreateIdea(idea *models.Idea) (*models.Idea, error) {
	result, err := ic.ideacollection.InsertOne(ic.ctx, idea)
	if err != nil {
		return nil, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to fetch inserted Idea _id")
	}
	idea.ID = oid
	return idea, err
}
