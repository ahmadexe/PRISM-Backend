package models

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	UserId       string             `json:"userId" bson:"userId" validate:"required"`
	ImageUrl     *string            `json:"imageUrl" bson:"imageUrl" validate:"required_without=Caption"`
	Caption      *string            `json:"caption" bson:"caption" validate:"required_without=ImageUrl"`
	NoOfViews    int                `json:"noOfViews" bson:"noOfViews"`
	UpVotes      int                `json:"noOfLikes" bson:"noOfLikes"`
	DownVotes    int                `json:"noOfDislikes" bson:"noOfDislikes"`
	NoOfComments int                `json:"noOfComments" bson:"noOfComments"`
	IsBanned     bool               `json:"isBanned" bson:"isBanned"`
	TotalReports int                `json:"totalReports" bson:"totalReports"`
	CreatedAt    int64              `json:"createdAt" bson:"createdAt"`
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
