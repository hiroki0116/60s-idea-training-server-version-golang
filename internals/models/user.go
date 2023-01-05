package models

import (
	"idea-training-version-go/internals/guard"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirebaseUID string             `json:"firebaseUid,omitempty" bson:"firebaseUid,omitempty"`
	FirstName   string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName    string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email       string             `json:"email,omitempty" validate:"required,email" bson:"email,omitempty"`
	Role        guard.Role         `json:"role,omitempty" bson:"role,omitempty"`
	Images      []Image            `json:"images" bson:"images"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type Image struct {
	About string `json:"about" bson:"about"`
	Url   string `json:"url" bson:"url"`
}

func (u *User) MarshalBSON() ([]byte, error) {

	u.UpdatedAt = time.Now()

	type custom User
	return bson.Marshal((*custom)(u))
}
