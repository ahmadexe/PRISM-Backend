package repository

import "go.mongodb.org/mongo-driver/mongo"

type ChatRepo struct {
	conversationsCol    *mongo.Collection
	messagesCol *mongo.Collection
}

func InitChatRepo(client *mongo.Client) *ChatRepo {
	conversationsCol := client.Database("chat-db").Collection("conversations")
	messagesCol := client.Database("chat-db").Collection("messages")
	return &ChatRepo{conversationsCol: conversationsCol, messagesCol: messagesCol}
}