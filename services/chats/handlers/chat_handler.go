package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ahmadexe/prism-backend/services/chats/data"
	"github.com/ahmadexe/prism-backend/services/chats/repository"
	"github.com/ahmadexe/prism-backend/services/chats/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	repo    *repository.ChatRepo
	rdb     *redis.Client
	clients map[string]*websocket.Conn
}

var broadcast chan data.Message

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitChatHandler(repo *repository.ChatRepo) *ChatHandler {

	broadcast = make(chan data.Message)

	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	clients := make(map[string]*websocket.Conn)
	return &ChatHandler{repo: repo, rdb: rdb, clients: clients}
}

func (handler *ChatHandler) HandleConnections(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	params := r.URL.Query()
	id1 := params.Get("id1")
	handler.clients[id1] = conn

	id2 := params.Get("id2")

	primitiveId1, err := primitive.ObjectIDFromHex(id1)
	if err != nil {
		log.Println(err)
		return
	}

	primitiveId2, err := primitive.ObjectIDFromHex(id2)
	if err != nil {
		log.Println(err)
		return
	}

	id := utils.SortIDs(primitiveId1, primitiveId2)
	var msgs []data.Message
	fetch, _ := handler.rdb.Get(ctx, id).Result()

	err = json.Unmarshal([]byte(fetch), &msgs)
	if err != nil {
		log.Println("No messages found")
	} else {
		handler.repo.PushBulkMessages(ctx, msgs)

		handler.rdb.Del(ctx, id)
		delete(handler.clients, id)
	}

	for {
		var msg data.Message
		err := conn.ReadJSON(&msg)
		if err != nil {

			fetch, _ := handler.rdb.Get(ctx, id).Result()
			log.Printf("Disconnected by: %s", id)
			var messages []data.Message
			err = json.Unmarshal([]byte(fetch), &messages)
			if err != nil {
				log.Println("No messages found")
			} else {
				handler.repo.PushBulkMessages(ctx, messages)

				handler.rdb.Del(ctx, id)
				delete(handler.clients, id)
			}

			return
		}

		msg.Id = primitive.NewObjectID()
		er := msg.Validate()
		if er != nil {
			log.Println(er)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
			continue
		}
		broadcast <- msg
	}
}

func (handler *ChatHandler) HandleMessages() {
	c := context.Background()
	for {
		msg := <-broadcast
		receiverClient := handler.clients[msg.ReceiverId.Hex()]
		senderClient := handler.clients[msg.SenderId.Hex()]

		if receiverClient == nil && senderClient == nil {
			continue
		} else {
			if senderClient != nil {
				err := senderClient.WriteJSON(msg)
				if err != nil {
					log.Println(err)
					delete(handler.clients, msg.SenderId.Hex())
				}
			}

			if receiverClient != nil {
				e := receiverClient.WriteJSON(msg)
				if e != nil {
					log.Println(e)
					delete(handler.clients, msg.ReceiverId.Hex())
				}
			}

			var allMessages []data.Message

			ids := utils.SortIDs(msg.SenderId, msg.ReceiverId)

			prevMessages, er := handler.rdb.Get(c, ids).Result()
			if er != nil {
				allMessages = append(allMessages, msg)
			} else {
				err := json.Unmarshal([]byte(prevMessages), &allMessages)
				if err != nil {
					log.Println(err)
				}
				allMessages = append(allMessages, msg)
			}
			messagesJson, err := json.Marshal(allMessages)
			if err != nil {
				log.Println(err)
			}

			handler.rdb.Set(c, ids, messagesJson, 0)
		}
	}
}

func (handler *ChatHandler) HandleConversation(ctx *gin.Context) {
	var convo data.Conversation

	if err := ctx.ShouldBindJSON(&convo); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	if err := convo.Validate(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	handler.repo.CreateOrFetchConversation(ctx, convo)
}
