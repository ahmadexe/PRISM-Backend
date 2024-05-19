package repository

import (
	"context"
	"time"

	"github.com/ahmadexe/prism-backend/services/chats/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepo struct {
	conversationsCol *mongo.Collection
	messagesCol      *mongo.Collection
}

func InitChatRepo(client *mongo.Client) *ChatRepo {
	conversationsCol := client.Database("chat-db").Collection("conversations")
	messagesCol := client.Database("chat-db").Collection("messages")
	return &ChatRepo{conversationsCol: conversationsCol, messagesCol: messagesCol}
}

func (chatRepo *ChatRepo) GetConversations() {

}

func (chatRepo *ChatRepo) CreateOrFetchConversation(ctx *gin.Context, convo data.Conversation) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user1id := convo.User1Id
	user2id := convo.User2Id

	filter := bson.M{"user1Id": user1id, "user2Id": user2id}

	var existingConvo data.Conversation
	err := chatRepo.conversationsCol.FindOne(c, filter).Decode(&existingConvo)
	if err == mongo.ErrNoDocuments {
		chatRepo.conversationsCol.InsertOne(c, convo)
		ctx.JSON(200, gin.H{"message": "Conversation created successfully.", "new": true, "conversation": convo})
	} else {
		convoId := existingConvo.Id
		var messages []data.Message
		filter := bson.M{"conversationId": convoId}

		cursor, err := chatRepo.messagesCol.Find(c, filter)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Error fetching messages."})
			return
		}

		if err = cursor.All(c, &messages); err != nil {
			ctx.JSON(500, gin.H{"error": "Error fetching messages."})
			return
		}

		ctx.JSON(200, gin.H{"conversation": existingConvo, "messages": messages})
	}
}

func (chatRepo *ChatRepo) PushBulkMessages(ctx *gin.Context, messages []data.Message) {
	c := context.Background()

	var docs []interface{}
	for _, message := range messages {
		docs = append(docs, message)
	}

	// go func() {
	// 	_, err := chatRepo.messagesCol.InsertMany(c, docs)
	// 	if err != nil {
	// 		ctx.JSON(500, gin.H{"error": "Error adding messages."})
	// 		return
	// 	}
	// } ()

	_, err := chatRepo.messagesCol.InsertMany(c, docs)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error adding messages."})
		return
	}

	ctx.JSON(200, gin.H{"message": "Messages added successfully."})
}