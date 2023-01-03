package controllers

import (
	"context"
	"errors"
	"idea-training-version-go/internals/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdeaController struct {
	ideacollection *mongo.Collection
	ctx            context.Context
}

type IIdeaController interface {
	CreateIdea(idea *models.Idea) (*models.Idea, error)
	GetAllIdeas(userID primitive.ObjectID) ([]*models.Idea, error)
	GetIdeaByID(ideaID primitive.ObjectID) (*models.Idea, error)
	UpdateIdea(idea *models.Idea) error
	DeleteIdea(ideaID primitive.ObjectID) error
}

func NewIdeaController(ideacollection *mongo.Collection, ctx context.Context) IIdeaController {
	return &IdeaController{
		ideacollection: ideacollection,
		ctx:            ctx,
	}
}

func (ic *IdeaController) CreateIdea(idea *models.Idea) (*models.Idea, error) {
	idea.CreatedAt = time.Now()
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

func (ic *IdeaController) GetAllIdeas(userID primitive.ObjectID) ([]*models.Idea, error) {
	var ideas []*models.Idea

	query := bson.D{
		bson.E{
			Key:   "createdBy",
			Value: userID,
		},
	}

	count, _ := ic.ideacollection.CountDocuments(ic.ctx, query)
	if count == 0 {
		// return empty array if no tasks found
		return []*models.Idea{}, nil
	}

	cursor, err := ic.ideacollection.Find(ic.ctx, query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ic.ctx, &ideas); err != nil {
		return nil, err
	}

	return ideas, nil
}

func (ic *IdeaController) GetIdeaByID(ideaID primitive.ObjectID) (*models.Idea, error) {
	var idea models.Idea

	query := bson.D{
		bson.E{
			Key:   "_id",
			Value: ideaID,
		},
	}

	err := ic.ideacollection.FindOne(ic.ctx, query).Decode(&idea)
	if err != nil {
		return nil, err
	}

	return &idea, nil
}

func (ic *IdeaController) UpdateIdea(idea *models.Idea) error {
	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: idea.ID,
		},
	}

	_, err := ic.ideacollection.UpdateOne(ic.ctx, filter, bson.M{"$set": idea})
	return err
}

func (ic *IdeaController) DeleteIdea(ideaID primitive.ObjectID) error {
	filter := bson.D{
		bson.E{
			Key:   "_id",
			Value: ideaID,
		},
	}

	_, err := ic.ideacollection.DeleteOne(ic.ctx, filter)
	return err
}
