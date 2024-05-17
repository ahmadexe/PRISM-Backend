package routes

import (
	"net/http"

	"github.com/ahmadexe/prism-backend/services/chats/handlers"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	chatHandler *handlers.ChatHandler
	router 	   *gin.Engine
}

func InitChatRouter(chatHandler *handlers.ChatHandler, r *gin.Engine) *ChatRouter {
	return &ChatRouter{chatHandler: chatHandler, router: r}
}

func (r *ChatRouter) SetupRoutes() {
	chats := r.router.Group("/v1")
	{
		chats.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Welcome to Prism Chat Service",
			})
		})
		chats.GET("/ws", r.chatHandler.HandleConnections)
	}
}