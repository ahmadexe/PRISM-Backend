package models

import (
	"encoding/json"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	Id             primitive.ObjectID  `json:"id" bson:"_id"`
	PostId         string  `json:"postId" bson:"postId" validate:"required"`
	UserId         string  `json:"userId" bson:"userId" validate:"required"`
	Content        string  `json:"content" bson:"content" validate:"required"`
	CreatedAt      string  `json:"createdAt" bson:"createdAt" validate:"required"`
	UserProfilePic *string `json:"userProfilePic" bson:"userProfilePic"`
	UserName       string  `json:"userName" bson:"userName" validate:"required"`
}

func (comment *Comment) MarshalJSON() ([]byte, error) {
	type Alias Comment
	return json.Marshal(&struct {
		Id string `json:"id"`
		*Alias
	}{
		Id:    comment.Id.Hex(),
		Alias: (*Alias)(comment),
	})
}

func (comment *Comment) Validate() error {
	v := validator.New()
	return v.Struct(comment)
}
