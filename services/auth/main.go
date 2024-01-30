package main

import (
	"context"

	"github.com/ahmadexe/prism-backend/services/auth/configs"
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/ahmadexe/prism-backend/services/auth/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	configs := configs.InitConfigs()
	client := configs.SetupDB()
	authRepo := repositories.InitAuthRepo(client)
	authHanler := handlers.InitAuthHandler(authRepo)

	router := gin.Default()

	authRouter := routes.InitAuthRouter(authHanler, router)
	authRouter.SetupRoutes()

	router.Run(configs.Host + ":" + configs.Port)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
