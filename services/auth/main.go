package main

import (
	"context"
	// "fmt"

	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/ahmadexe/prism-backend/services/auth/routes"
	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://muahmad710:RiX91Y94aK9xVQW3@auth.dgg1cj0.mongodb.net/?retryWrites=true&w=majority"))

	if err != nil {
		panic(err)
	}

	authRepo :=  repositories.InitAuthRepo(client)
	authHanler := handlers.InitAuthHandler(authRepo)

	router := gin.Default()

	authRouter := routes.InitAuthRouter(authHanler, router)
	authRouter.SetupRoutes()

	router.Run(":8080")

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()
	// if err := client.Database("auth-service").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
