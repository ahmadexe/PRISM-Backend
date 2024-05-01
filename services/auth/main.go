package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/ahmadexe/prism-backend/services/auth/configs"
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/ahmadexe/prism-backend/services/auth/routes"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("../../env_var/app_keys.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	configs := configs.InitConfigs()
	client := configs.SetupDB()
	authRepo := repositories.InitAuthRepo(client)
	authHanler := handlers.InitAuthHandler(authRepo)

	router := gin.Default()

	authRouter := routes.InitAuthRouter(authHanler, router)
	authRouter.SetupRoutes(app)
	
	router.Run(configs.Host + ":" + configs.Port)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
