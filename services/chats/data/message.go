package data

import (
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt      int                `json:"createdAt" bson:"createdAt" validate:"required"`
	Message        string             `json:"message" bson:"message" validate:"required"`
	SenderId       primitive.ObjectID `json:"senderId" bson:"senderId" validate:"required"`
	ConversationId primitive.ObjectID `json:"conversationId" bson:"conversationId" validate:"required"`
	ReceiverId     primitive.ObjectID `json:"receiverId" bson:"receiverId" validate:"required"`
	SenderName     string             `json:"senderName" bson:"senderName" validate:"required"`
	ReceiverName   string             `json:"receiverName" bson:"receiverName" validate:"required"`
	SenderPic      *string            `json:"senderPic" bson:"senderPic"`
	ReceiverPic    *string            `json:"receiverPic" bson:"receiverPic"`
}

func (message *Message) Validate() error {
	v := validator.New()
	return v.Struct(message)
}
