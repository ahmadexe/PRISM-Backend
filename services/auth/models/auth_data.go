package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthData struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Fullname  string             `json:"fullname" bson:"fullname"`
	Domain    string             `json:"domain" bson:"domain"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
}

func (authData *AuthData) UnmarshalJSON([]byte) error  {
	// check if createdAt is not null
	// if null, set createdAt to current time
	// else, set createdAt to the value of createdAt

	if authData.CreatedAt == 0 {
		authData.CreatedAt = time.Now().UnixMicro()
	}
	return nil
}