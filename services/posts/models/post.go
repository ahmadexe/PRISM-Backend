package models

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	UserId       string             `json:"userId" bson:"userId" validate:"required"`
	ImageUrl     *string            `json:"imageUrl" bson:"imageUrl"`
	Caption      *string            `json:"caption" bson:"caption"`
	NoOfViews    int                `json:"noOfViews" bson:"noOfViews"`
	NoOfLikes    int                `json:"noOfLikes" bson:"noOfLikes"`
	NoOfComments int                `json:"noOfComments" bson:"noOfComments"`
	IsBanned     bool               `json:"isBanned" bson:"isBanned"`
	TotalReports int                `json:"totalReports" bson:"totalReports"`
}

func (post *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	return json.Marshal(&struct {
		Id string `json:"id"`
		*Alias
	}{
		Id:    post.Id.Hex(),
		Alias: (*Alias)(post),
	})
}

func (post *Post) Validate() error {
	v := validator.New()
	return v.Struct(post)
}