package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ahmadexe/prism-backend/services/chats/data"
	"github.com/ahmadexe/prism-backend/services/chats/repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type ChatHandler struct {
	repo *repository.ChatRepo
	rdb  *redis.Client
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
	id := params.Get("id")
	handler.clients[id] = conn

	for {
		var msg data.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fetch, _ := handler.rdb.Get(ctx, id).Result()

			var message []data.Message
			err = json.Unmarshal([]byte(fetch), &message)
			if err != nil {
				log.Println("No messages found")
			}

			fmt.Println(message)
			
			handler.rdb.Del(ctx, id)
			delete(handler.clients, id)
			return
		}

		broadcast <- msg
	}
}

func (handler *ChatHandler) HandleMessages() {
	c := context.Background()
	for {
		msg := <-broadcast

		client := handler.clients[msg.ReceiverId]
		if client == nil {
			continue
		} else {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				delete(handler.clients, msg.ReceiverId)
			}

			var allMessages []data.Message

			prevMessages, er := handler.rdb.Get(c, msg.SenderId).Result()
			if er != nil {
				allMessages = append(allMessages, msg)
			} else {
				err = json.Unmarshal([]byte(prevMessages), &allMessages)
				if err != nil {
					log.Println(err)
				}
				allMessages = append(allMessages, msg)
			}
			messagesJson, err := json.Marshal(allMessages)
			if err != nil {
				log.Println(err)
			}

			handler.rdb.Set(c, msg.SenderId, messagesJson, 0)
		}
	}
}
