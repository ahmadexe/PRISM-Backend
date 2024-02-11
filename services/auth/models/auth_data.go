package models

import (
	"encoding/json"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthData struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	Uid            *string            `json:"uid" bson:"uid" validate:"required"`
	Email          *string            `json:"email" bson:"email" validate:"required"`
	Fullname       *string            `json:"fullname" validate:"required" bson:"fullname"`
	Domain         *string            `json:"domain" validate:"required" bson:"domain"`
	Bio            *string             `json:"bio" bson:"bio"`
	ImageUrl       *string            `json:"imageUrl" bson:"imageUrl"`
	BannerImageUrl *string            `json:"bannerImageUrl" bson:"bannerImageUrl"`
	Followers      []string           `json:"followers" bson:"followers"`
	Following      []string           `json:"following" bson:"following"`
	CreatedAt      int64              `json:"createdAt" bson:"createdAt"`
}

func (authData *AuthData) MarshalJSON() ([]byte, error) {
	type Alias AuthData
	return json.Marshal(&struct {
		Id string `json:"id"`
		*Alias
	}{
		Id:    authData.Id.Hex(),
		Alias: (*Alias)(authData),
	})
}

func (authData *AuthData) Validate() error {
	v := validator.New()
	return v.Struct(authData)
}
