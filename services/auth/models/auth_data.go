package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthData struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email string `json:"email" bson:"email"`
	Fullname string `json:"fullname" bson:"fullname"`
	Domain string `json:"domain" bson:"domain"`

}