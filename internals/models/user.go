package models

import (
	"idea-training-version-go/internals/guard"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const DEFAULT_USER_ROLE = "user"
const DEFAULT_USER_IMAGE = "https://res.cloudinary.com/sixty-seconds-idea-training-project/image/upload/v1656157889/users/default-user-image_LYizIFTei_ioicfh.png"

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirebaseUID string             `json:"firebaseUid" bson:"firebaseUid"`
	FirstName   string             `json:"firstName" validate:"required" bson:"firstName"`
	LastName    string             `json:"lastName" validate:"required" bson:"lastName"`
	Email       string             `json:"email" validate:"required,email" bson:"email"`
	Role        guard.Role         `json:"role" bson:"role"`
	Images      []Image            `json:"images" bson:"images"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type Image struct {
	About string `json:"about" bson:"about"`
	Url   string `json:"url" bson:"url"`
}

func (u *User) MarshalBSON() ([]byte, error) {

	// deal with time stamps
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	// deal with default role
	u.Role = DEFAULT_USER_ROLE

	// deal with default image
	if len(u.Images) == 0 {
		u.Images = append(u.Images, Image{
			About: "default",
			Url:   DEFAULT_USER_IMAGE,
		})
	}

	type custom User
	return bson.Marshal((*custom)(u))
}
