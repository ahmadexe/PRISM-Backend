package main

import (
	"context"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"github.com/ahmadexe/prism-backend/services/auth/configs"
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/ahmadexe/prism-backend/services/auth/routes"
	"github.com/gin-contrib/cors"
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
	searchHandler := handlers.InitSearchHandler(authRepo)
	authHanler := handlers.InitAuthHandler(authRepo)

	go searchHandler.HandleSearch()
	gin.SetMode(configs.Mode)
	router := gin.Default()

	authRouter := routes.InitAuthRouter(authHanler, searchHandler, router)
	authRouter.SetupRoutes(app)
	
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, 
	}

	router.Use(cors.New(corsConfig))

	router.Run(configs.Host + ":" + configs.Port)

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
}
