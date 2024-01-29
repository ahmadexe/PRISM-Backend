package repositories

import (
	"log"

	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	Collection *mongo.Collection
}

func InitAuthRepo(client *mongo.Client) *AuthRepo {
	collection := client.Database("auth-service").Collection("users")
	return &AuthRepo{Collection: collection}
}

func (repo *AuthRepo) AddUser(ctx *gin.Context) {
	var user models.AuthData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Printf("%+v", user)

	if err := user.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := repo.Collection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User added successfully", "data": result})
}
