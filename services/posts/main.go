package main

import (
	"time"

	"github.com/ahmadexe/prism-backend/services/posts/configs"
	"github.com/ahmadexe/prism-backend/services/posts/handlers"
	"github.com/ahmadexe/prism-backend/services/posts/repositories"
	"github.com/ahmadexe/prism-backend/services/posts/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configs := configs.InitConfigs()
	client := configs.SetupDB()

	postRepo := repositories.InitPostRepo(client)
	postHanlder := handlers.InitPostHandler(postRepo)

	commentHandler := handlers.InitCommentHandler(postRepo)

	gin.SetMode(configs.Mode)
	router := gin.Default()

	postRouter := routes.InitPostsRouter(postHanlder, router)
	postRouter.SetupRoutes()

	commentRouter := routes.InitCommentRouter(commentHandler, router)
	commentRouter.SetupRoutes()
	
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
