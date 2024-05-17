package main

import (
	"github.com/ahmadexe/prism-backend/services/chats/configs"
	"github.com/ahmadexe/prism-backend/services/chats/handlers"
	"github.com/ahmadexe/prism-backend/services/chats/repository"
	"github.com/ahmadexe/prism-backend/services/chats/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	configs := configs.InitConfigs()
	client := configs.SetupDB()

	repo := repository.InitChatRepo(client)
	handler := handlers.InitChatHandler(repo)
	go handler.HandleMessages()
	router := gin.Default()

	chatRouter := routes.InitChatRouter(handler, router)
	chatRouter.SetupRoutes()

	router.Run(configs.Host + ":" + configs.Port)
}
