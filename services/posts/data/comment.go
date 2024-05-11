package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	PostId         primitive.ObjectID `json:"postId" bson:"postId" validate:"required"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
	Content        string             `json:"content" bson:"content" validate:"required"`
	CreatedAt      int                `json:"createdAt" bson:"createdAt" validate:"required"`
	UserProfilePic *string            `json:"userProfilePic" bson:"userProfilePic"`
	UserName       string             `json:"userName" bson:"userName" validate:"required"`
}

func (comment *Comment) Validate() error {
	v := validator.New()
	return v.Struct(comment)
}
