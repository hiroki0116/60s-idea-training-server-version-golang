package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Idea struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TopicTitle string             `json:"topicTitle,omitempty" bson:"topicTitle,omitempty"`
	Category   string             `json:"category,omitempty" bson:"category,omitempty"`
	Ideas      *[]string          `json:"ideas,omitempty" bson:"ideas,omitempty"`
	CreatedBy  primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Viewed     *bool              `json:"viewed,omitempty" bson:"viewed,omitempty"`
	IsLiked    *bool              `json:"isLiked,omitempty" bson:"isLiked,omitempty"`
	Comment    *string            `json:"comment,omitempty" bson:"comment,omitempty"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (i *Idea) MarshalBSON() ([]byte, error) {

	// deal with time stamps
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now()
	}
	i.UpdatedAt = time.Now()

	// deal with default topic title
	if i.TopicTitle == "" {
		i.TopicTitle = "Untitled"
	}

	// deal with default category
	if i.Category == "" {
		i.Category = "Other"
	}
	// deal with default viewed
	i.Viewed = &[]bool{false}[0]

	// deal with default isLiked
	i.IsLiked = &[]bool{false}[0]

	// deal with default ideas
	if i.Ideas == nil {
		i.Ideas = &[]string{}
	}

	// deal with default comment
	i.Comment = &[]string{""}[0]

	type custom Idea
	return bson.Marshal((*custom)(i))
}
