package models

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	UserId         string             `json:"userId" bson:"userId" validate:"required"`
	Title          string             `json:"title" bson:"title" validate:"required"`
	ImageUrl       *string            `json:"imageUrl" bson:"imageUrl" validate:"required_without=Description"`
	Description    *string            `json:"description" bson:"description" validate:"required_without=ImageUrl"`
	NoOfViews      int                `json:"noOfViews" bson:"noOfViews"`
	UserName       string             `json:"userName" bson:"userName" validate:"required"`
	Category       string             `json:"category" bson:"category" validate:"required"`
	UserProfilePic string             `json:"userProfilePic" bson:"userProfilePic"`
	UpVotes        int                `json:"upVotes" bson:"upVotes"`
	DownVotes      int                `json:"downVotes" bson:"downVotes"`
	UpVotedBy      []string           `json:"upVotedBy" bson:"upVotedBy"`
	DownVotedBy    []string           `json:"downVotedBy" bson:"downVotedBy"`
	NoOfComments   int                `json:"noOfComments" bson:"noOfComments"`
	CommentedBy    []string           `json:"commentedBy" bson:"commentedBy"`
	IsBanned       bool               `json:"isBanned" bson:"isBanned"`
	TotalReports   int                `json:"totalReports" bson:"totalReports"`
	CreatedAt      int64              `json:"createdAt" bson:"createdAt"`
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
