package main

import (
	"time"

	"github.com/ahmadexe/prism-backend/services/chats/configs"
	"github.com/ahmadexe/prism-backend/services/chats/handlers"
	"github.com/ahmadexe/prism-backend/services/chats/repository"
	"github.com/ahmadexe/prism-backend/services/chats/routes"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	configs := configs.InitConfigs()
	client := configs.SetupDB()

	repo := repository.InitChatRepo(client)
	handler := handlers.InitChatHandler(repo)
	go handler.HandleMessages()
	gin.SetMode(configs.Mode)
	router := gin.Default()

	chatRouter := routes.InitChatRouter(handler, router)
	chatRouter.SetupRoutes()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, 
	}

	router.Use(cors.New(corsConfig))

	router.Run(configs.Host + ":" + configs.Port)
}
