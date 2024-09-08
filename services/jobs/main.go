package main

import (
	"context"

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

	router := gin.Default()
	jobsRouter := routes.InitJobsRouter(handler, router)
	jobsRouter.SetupRoutes()

	router.Use(cors.Default())
	router.Run(configs.Host + ":" + configs.Port)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}