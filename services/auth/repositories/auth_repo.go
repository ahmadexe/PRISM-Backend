package repositories

import (
	"net/http"

	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	Collection *mongo.Collection
}

func InitAuthRepo(client *mongo.Client) *AuthRepo {
	collection := client.Database("prism-dev").Collection("users")
	return &AuthRepo{Collection: collection}
}

func (repo *AuthRepo) AddUser(user models.AuthData, ctx *gin.Context) {
	result, err := repo.Collection.InsertOne(ctx, user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User added successfully", "data": result})
}

func (repo *AuthRepo) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	filter := bson.M{"uid": id}
	var user models.AuthData
	
	if err := repo.Collection.FindOne(ctx, filter).Decode(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User found successfully", "data": user})
}

func (repo *AuthRepo) UpdateUser(user models.AuthData, ctx *gin.Context) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	if _, err := repo.Collection.UpdateOne(ctx, filter, update); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}