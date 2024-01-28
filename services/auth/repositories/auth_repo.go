package repositories

import (
	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	Client *mongo.Client
}

func InitAuthRepo(client *mongo.Client) *AuthRepo {
	return &AuthRepo{Client: client}
}

func (repo *AuthRepo) AddUser(ctx *gin.Context) {
	collection := repo.Client.Database("auth-service").Collection("users")
	var user models.AuthData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "User added successfully"})
}
