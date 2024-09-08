package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Request struct {
	ID     primitive.ObjectID `json:"id" bson:"_id" validate:"required"`
	UserId primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
}

func (app *Request) Validate() error {
	v := validator.New()
	err := v.Struct(app)
	return err
}
