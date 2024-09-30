package main

import (
	"context"
	"time"

	"github.com/ahmadexe/prism-backend/services/jobs/configs"
	"github.com/ahmadexe/prism-backend/services/jobs/handlers"
	"github.com/ahmadexe/prism-backend/services/jobs/repository"
	"github.com/ahmadexe/prism-backend/services/jobs/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main()  {
	configs := configs.InitConfigs()
	client := configs.SetupDB()
	
	jobsRepo := repository.NewJobsRepo(client)

	handler := handlers.InitJobHandler(jobsRepo)
	
	gin.SetMode(configs.Mode)
	router := gin.Default()
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))
	

	jobsRouter := routes.InitJobsRouter(handler, router)
	jobsRouter.SetupRoutes()

	router.Run(configs.Host + ":" + configs.Port)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}