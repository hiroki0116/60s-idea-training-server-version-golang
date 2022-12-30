package models

import (
	"idea-training-version-go/internals/guard"
	"time"
)

const DEFAULT_USER_ROLE = "user"
const DEFAULT_USER_IMAGE = "https://res.cloudinary.com/sixty-seconds-idea-training-project/image/upload/v1656157889/users/default-user-image_LYizIFTei_ioicfh.png"

type User struct {
	FirebaseUID string     `json:"firebaseUid" bson:"firebaseUid"`
	FirstName   string     `json:"firstName" validate:"required" bson:"firstName"`
	LastName    string     `json:"lastName" validate:"required" bson:"lastName"`
	Email       string     `json:"email" validate:"required,email" bson:"email"`
	Role        guard.Role `json:"role" bson:"role"`
	Images      []Image    `json:"images" bson:"images"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updatedAt"`
	CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
}

type Image struct {
	About string `json:"about" bson:"about"`
	Url   string `json:"url" bson:"url"`
}
