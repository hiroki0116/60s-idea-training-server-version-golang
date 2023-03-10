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
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func (i *Idea) MarshalBSON() ([]byte, error) {
	// deal with time stamps
	i.UpdatedAt = time.Now()

	type custom Idea
	return bson.Marshal((*custom)(i))
}
