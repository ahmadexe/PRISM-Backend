package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchRequest struct {
	Query string `json:"query" bson:"query" validate:"required"`
	Id primitive.ObjectID `json:"id" bson:"id" validate:"required"`
}

func (sr *SearchRequest) Validate() error {
	v := validator.New()
	return v.Struct(sr)
}