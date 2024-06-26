package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowRequest struct {
	From primitive.ObjectID `json:"from" bson:"from" validate:"required"`
	To   primitive.ObjectID `json:"to" bson:"to" validate:"required"`
}

func (fr *FollowRequest) Validate() error {
	v := validator.New()
	return v.Struct(fr)
}
