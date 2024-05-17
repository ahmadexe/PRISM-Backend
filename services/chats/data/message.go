package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt      int                `json:"createdAt" bson:"createdAt"`
	Message        string             `json:"message" bson:"message"`
	SenderId       string             `json:"senderId" bson:"senderId"`
	ConversationId string             `json:"conversationId" bson:"conversationId"`
	ReceiverId     string             `json:"receiverId" bson:"receiverId"`
}
