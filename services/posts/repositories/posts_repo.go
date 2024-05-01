package repositories

import (
	"context"
	"log"
	"net/http"
	"time"

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

func (repo *PostRepo) AddPost(post models.Post, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	_, err := repo.collection.InsertOne(c, post)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding post to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}