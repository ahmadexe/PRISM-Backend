package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenRequest struct {
	DeviceToken string             `json:"deviceToken" bson:"deviceToken" validate:"required"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId" validate:"required"`
}

func (tr *TokenRequest) Validate() error {
	v := validator.New()
	return v.Struct(tr)
}	
