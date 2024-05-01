package repositories

import (
	"log"
	"net/http"

	"github.com/ahmadexe/prism-backend/services/posts/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo struct {
	collection *mongo.Collection
}

func InitPostRepo(client *mongo.Client) *PostRepo {
	collection := client.Database("prism-dev").Collection("posts")
	return &PostRepo{collection: collection}
}

// AddPost adds a post to the database
func (repo *PostRepo) AddPost(post models.Post, ctx *gin.Context) {
	_, err := repo.collection.InsertOne(ctx, post)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding post to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}