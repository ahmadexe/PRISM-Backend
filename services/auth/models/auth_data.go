package models

import (
	"encoding/json"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthData struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Email     *string             `json:"email" bson:"email" validate:"required"`
	Fullname  *string             `json:"fullname" validate:"required" bson:"fullname"`
	Domain    *string             `json:"domain" validate:"required" bson:"domain"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
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

func (authData *AuthData) UnmarshalJSON(d []byte) error {
	// check if createdAt is not null
	// if null, set createdAt to current time
	// else, set createdAt to the value of createdAt

	type Alias AuthData
	ad := &struct{
		Id primitive.ObjectID `json:"id"`
		CreatedAt int64 `json:"createdAt"`
		*Alias
	}{
		CreatedAt: time.Now().UnixMicro(),
		Id: primitive.NewObjectID(),
		Alias: (*Alias)(authData),
	}

	if err := json.Unmarshal(d, &ad); err != nil {
		return err
	}

	authData.Id = ad.Id
	authData.CreatedAt = ad.CreatedAt

	return nil
}

func (authData *AuthData) Validate() error {
	v := validator.New()
	return v.Struct(authData)
}