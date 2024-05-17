package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	User1Id   primitive.ObjectID `json:"user1Id" bson:"user1Id" validate:"required"`
	User2Id   primitive.ObjectID `json:"user2Id" bson:"user2Id" validate:"required"`
	User1Name string             `json:"user1Name" bson:"user1Name" validate:"required"`
	User2Name string             `json:"user2Name" bson:"user2Name" validate:"required"`
	User1Pic  *string            `json:"user1Pic" bson:"user1Pic"`
	User2Pic  *string            `json:"user2Pic" bson:"user2Pic"`
}

func (conversation *Conversation) Validate() error {
	v := validator.New()
	return v.Struct(conversation)
}
