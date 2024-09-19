package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportRequest struct {
	PostId primitive.ObjectID `json:"postId" bson:"postId" validate:"required"`
	Type   string `json:"type" bson:"type" validate:"required"`
}

func (reportRequest *ReportRequest) Validate() error {
	v := validator.New()
	return v.Struct(reportRequest)
}
