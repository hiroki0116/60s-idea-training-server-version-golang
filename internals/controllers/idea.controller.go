package controllers

import (
	"context"
	"idea-training-version-go/internals/models"
	"log"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	errors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetTotalIdeasOfToday(userID primitive.ObjectID) ([]bson.M, error)
	GetTotalIdeasOfAllTime(userID primitive.ObjectID) ([]bson.M, error)
	GetTotalConsecutiveDays(userID primitive.ObjectID) (int, error)
	GetRecentIdeas(userID primitive.ObjectID) ([]*models.Idea, error)
	GetWeeklyIdeas(userID primitive.ObjectID) ([]bson.M, time.Time, error)
	Search(filter bson.M, sort int, page int, limit int) ([]models.Idea, *paginate.PaginatedData, error)
}

func NewIdeaController(ideacollection *mongo.Collection, ctx context.Context) IIdeaController {
	return &IdeaController{
		ideacollection: ideacollection,
		ctx:            ctx,
	}
}

func (ic *IdeaController) CreateIdea(idea *models.Idea) (*models.Idea, error) {
	idea.CreatedAt = time.Now()
	// deal with default topic title
	if idea.TopicTitle == "" {
		idea.TopicTitle = "Untitled"
	}

	// deal with default category
	if idea.Category == "" {
		idea.Category = "Other"
	}
	// deal with default viewed
	if idea.Viewed == nil {
		idea.Viewed = &[]bool{false}[0]
	}

	// deal with default isLiked
	if idea.IsLiked == nil {
		idea.IsLiked = &[]bool{false}[0]
	}

	// deal with default ideas
	if idea.Ideas == nil {
		idea.Ideas = &[]string{}
	}

	// deal with default comment
	if idea.Comment == nil {
		idea.Comment = &[]string{""}[0]
	}

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

func (ic *IdeaController) GetTotalIdeasOfToday(userID primitive.ObjectID) ([]bson.M, error) {
	now := time.Now()
	year, month, day := now.UTC().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	matchStage := bson.D{
		bson.E{
			Key: "$match",
			Value: bson.D{
				bson.E{
					Key:   "createdBy",
					Value: userID,
				},
				bson.E{
					Key: "createdAt",
					Value: bson.D{
						bson.E{
							Key:   "$gte",
							Value: today,
						},
					},
				},
			},
		},
	}

	projectgStage := bson.D{
		bson.E{
			Key: "$group",
			Value: bson.D{
				bson.E{
					Key:   "_id",
					Value: nil,
				},
				bson.E{
					Key: "totalIdeas",
					Value: bson.D{
						bson.E{
							Key: "$sum",
							Value: bson.D{
								bson.E{
									Key:   "$size",
									Value: "$ideas",
								},
							},
						},
					},
				},
				bson.E{
					Key: "totalSessions",
					Value: bson.D{
						bson.E{
							Key:   "$sum",
							Value: 1,
						},
					},
				},
			},
		},
	}

	cursor, err := ic.ideacollection.Aggregate(ic.ctx, mongo.Pipeline{matchStage, projectgStage})
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}

	return results, nil
}

func (ic *IdeaController) GetTotalIdeasOfAllTime(userID primitive.ObjectID) ([]bson.M, error) {
	matchStage := bson.D{
		bson.E{
			Key: "$match",
			Value: bson.D{
				bson.E{
					Key:   "createdBy",
					Value: userID,
				},
			},
		},
	}

	projectgStage := bson.D{
		bson.E{
			Key: "$group",
			Value: bson.D{
				bson.E{
					Key:   "_id",
					Value: nil,
				},
				bson.E{
					Key: "totalIdeas",
					Value: bson.D{
						bson.E{
							Key: "$sum",
							Value: bson.D{
								bson.E{
									Key:   "$size",
									Value: "$ideas",
								},
							},
						},
					},
				},
				bson.E{
					Key: "totalSessions",
					Value: bson.D{
						bson.E{
							Key:   "$sum",
							Value: 1,
						},
					},
				},
			},
		},
	}

	cursor, err := ic.ideacollection.Aggregate(ic.ctx, mongo.Pipeline{matchStage, projectgStage})
	if err != nil {
		return nil, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, err
	}

	return results, nil
}

func (ic *IdeaController) GetTotalConsecutiveDays(userID primitive.ObjectID) (int, error) {
	now := time.Now()
	year, month, day := now.UTC().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	var consecutiveDays int = 0
	var isConsecutive bool = true

	for {
		dateFrom := today.AddDate(0, 0, (-1)*consecutiveDays-1)
		dateTo := today.AddDate(0, 0, (-1)*consecutiveDays)

		filter := bson.D{
			{
				Key:   "createdBy",
				Value: userID,
			},
			{
				Key: "createdAt",
				Value: bson.D{
					{
						Key:   "$gte",
						Value: dateFrom,
					},
					{
						Key:   "$lt",
						Value: dateTo,
					},
				},
			},
		}

		numOfDoc, err := ic.ideacollection.CountDocuments(ic.ctx, filter)

		if err != nil {
			return 0, errors.Wrap(err, "error while counting documents")
		}

		if numOfDoc > 0 {
			consecutiveDays++
		} else {
			isConsecutive = false
		}

		if !isConsecutive {
			break
		}
	}

	return consecutiveDays, nil
}

func (ic *IdeaController) GetRecentIdeas(userID primitive.ObjectID) ([]*models.Idea, error) {
	var ideas []*models.Idea

	query := bson.D{
		bson.E{
			Key:   "createdBy",
			Value: userID,
		},
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "createdAt", Value: -1}}).SetLimit(5)

	count, _ := ic.ideacollection.CountDocuments(ic.ctx, query)
	if count == 0 {
		// return empty array if no tasks found
		return []*models.Idea{}, nil
	}

	cursor, err := ic.ideacollection.Find(ic.ctx, query, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ic.ctx, &ideas); err != nil {
		return nil, err
	}

	return ideas, nil
}

func (ic *IdeaController) GetWeeklyIdeas(userID primitive.ObjectID) ([]bson.M, time.Time, error) {
	// get first day of week
	weekday := time.Duration(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := time.Now().Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	lastMonday := currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)

	matchStage := bson.D{
		bson.E{
			Key: "$match",
			Value: bson.D{
				bson.E{
					Key:   "createdBy",
					Value: userID,
				},
				bson.E{
					Key: "createdAt",
					Value: bson.D{
						bson.E{
							Key:   "$gte",
							Value: lastMonday,
						},
					},
				},
			},
		},
	}
	groupStage := bson.D{
		bson.E{
			Key: "$group",
			Value: bson.D{
				bson.E{
					Key: "_id",
					Value: bson.D{
						bson.E{
							Key: "$dateToString",
							Value: bson.D{
								bson.E{
									Key:   "format",
									Value: "%Y-%m-%d",
								},
								bson.E{
									Key:   "date",
									Value: "$createdAt",
								},
							},
						},
					},
				},
				bson.E{
					Key: "totalIdeas",
					Value: bson.D{
						bson.E{
							Key: "$sum",
							Value: bson.D{
								bson.E{
									Key:   "$size",
									Value: "$ideas",
								},
							},
						},
					},
				},
				bson.E{
					Key: "totalSessions",
					Value: bson.D{
						bson.E{
							Key:   "$sum",
							Value: 1,
						},
					},
				},
			},
		},
	}

	cursor, err := ic.ideacollection.Aggregate(ic.ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, lastMonday, err
	}

	// display the results
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return results, lastMonday, err
	}

	return results, lastMonday, nil

}

func (ic *IdeaController) Search(filter bson.M, sort int, page int, limit int) ([]models.Idea, *paginate.PaginatedData, error) {

	collation := options.Collation{
		Locale:   "en",
		Strength: 1,
	}
	var ideas []models.Idea
	paginatedData, err := paginate.New(ic.ideacollection).SetCollation(&collation).Context(ic.ctx).Limit(int64(limit)).Page(int64(page)).Sort("updatedAt", sort).Filter(filter).Decode(&ideas).Find()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	return ideas, paginatedData, nil
}
