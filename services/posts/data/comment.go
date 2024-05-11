package data

import (
	// "encoding/json"

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

// func (comment *Comment) MarshalJSON() ([]byte, error) {
// 	type Alias Comment
// 	return json.Marshal(&struct {
// 		Id string `json:"id"`
// 		*Alias
// 	}{
// 		Id:    comment.Id.Hex(),
// 		Alias: (*Alias)(comment),
// 	})
// }

// func (comment *Comment) UnmarshalJSON(data []byte) error {
// 	type Alias Comment
// 	aux := &struct {
// 		PostId string `json:"postId"`
// 		UserId string `json:"userId"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(comment),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	comment.PostId, _ = primitive.ObjectIDFromHex(aux.PostId)
// 	comment.UserId, _ = primitive.ObjectIDFromHex(aux.UserId)
// 	return nil
// }

func (comment *Comment) Validate() error {
	v := validator.New()
	return v.Struct(comment)
}
