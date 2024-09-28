package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobApplication struct {
	Id       primitive.ObjectID  `json:"id" bson:"_id"`
	JobId    string  `json:"jobId" bson:"jobId" validate:"required"`
	UserId   string  `json:"userId" bson:"userId" validate:"required"`
	IsHired  bool    `json:"isHired" bson:"isHired"`
	Username string  `json:"username" bson:"username" validate:"required"`
	Avatar   *string `json:"avatar" bson:"avatar"`
}

func (ja *JobApplication) Validate() error {
	v := validator.New()
	return v.Struct(ja)
}
