package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ahmadexe/prism-backend/services/auth/data"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/websocket"
)

type SearchHandler struct {
	repo    *repositories.AuthRepo
	rdb     *redis.Client
	clients map[string]*websocket.Conn
}

var broadcast chan data.SearchRequest

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitSearchHandler(repo *repositories.AuthRepo) *SearchHandler {
	broadcast = make(chan data.SearchRequest)

	rdb := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	clients := make(map[string]*websocket.Conn)

	return &SearchHandler{repo: repo, rdb: rdb, clients: clients}
}

func (handler *SearchHandler) HandleConnections(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	id := ctx.Param("id")
	fmt.Println("HERE 6")
	handler.clients[id] = conn
	fmt.Println("HERE 7")

	for {
		var msg data.SearchRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(handler.clients, id)
			break
		}

		broadcast <- msg
	}
}

func (handler *SearchHandler) HandleSearch() {
	for {
		msg := <-broadcast
		client := handler.clients[msg.Id.Hex()]

		if client == nil {
			continue
		}

		users, err := handler.repo.GetUserBySubString(msg.Query)
		if err != nil {
			log.Println(err)
			continue
		}

		response := map[string]interface{}{
			"data": users,
		}
		err = client.WriteJSON(response)
		if err != nil {
			log.Println(err)
			client.Close()
			delete(handler.clients, msg.Id.Hex())
		}
	}
}
