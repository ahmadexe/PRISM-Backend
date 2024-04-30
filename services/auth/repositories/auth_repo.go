package repositories

import (
	"log"
	"net/http"

	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	collection *mongo.Collection
}

// func InitAuthRepo(client *mongo.Client) *AuthRepo {
// 	collection := client.Database("prism-dev").Collection("users")
// 	return &AuthRepo{Collection: collection}
// }

func InitAuthRepo(client *mongo.Client) *AuthRepo {
	collection := client.Database("prism-dev").Collection("users")
	return &AuthRepo{collection: collection}
}

func (repo *AuthRepo) AddUser(user models.AuthData, ctx *gin.Context) {
	_, err := repo.collection.InsertOne(ctx, user)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (repo *AuthRepo) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	filter := bson.M{"uid": id}
	var user models.AuthData

	if err := repo.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (repo *AuthRepo) UpdateUser(user models.AuthData, ctx *gin.Context) {
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": user}

	if _, err := repo.collection.UpdateOne(ctx, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding user to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
