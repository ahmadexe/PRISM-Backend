package main

import (
	"github.com/ahmadexe/prism-backend/services/posts/configs"
	"github.com/ahmadexe/prism-backend/services/posts/handlers"
	"github.com/ahmadexe/prism-backend/services/posts/repositories"
	"github.com/ahmadexe/prism-backend/services/posts/routes"
	"github.com/gin-gonic/gin"
)

func main()  {
	configs := configs.InitConfigs()
	client := configs.SetupDB()

	postRepo := repositories.InitPostRepo(client)
	postHanlder := handlers.InitPostHandler(postRepo)

	router := gin.Default()

	postRouter := routes.InitPostsRouter(postHanlder, router)
	postRouter.SetupRoutes()

	router.Run(configs.Host + ":" + configs.Port)
}