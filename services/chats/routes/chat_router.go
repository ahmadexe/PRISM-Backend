package routes

import (
	"github.com/ahmadexe/prism-backend/services/chats/handlers"
	"github.com/ahmadexe/prism-backend/services/chats/middlewares"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	chatHandler *handlers.ChatHandler
	router      *gin.Engine
}

func InitChatRouter(chatHandler *handlers.ChatHandler, r *gin.Engine) *ChatRouter {
	return &ChatRouter{chatHandler: chatHandler, router: r}
}

func (r *ChatRouter) SetupRoutes() {
	chats := r.router.Group("/v1")
	{
		chats.GET("/chats/ws", r.chatHandler.HandleConnections)
		chats.POST("/chats/convo", middlewares.VerifyUser, r.chatHandler.HandleConversation)
		chats.GET("/chats/convo/:id", middlewares.VerifyUser, r.chatHandler.GetConversations)
	}
}
