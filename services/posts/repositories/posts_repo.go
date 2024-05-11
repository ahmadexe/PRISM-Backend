package repositories

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ahmadexe/prism-backend/services/posts/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo struct {
	collection *mongo.Collection
}

func InitPostRepo(client *mongo.Client) *PostRepo {
	collection := client.Database("posts-db").Collection("posts")
	return &PostRepo{collection: collection}
}

func (repo *PostRepo) AddPost(post models.Post, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := repo.collection.InsertOne(c, post)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding post to database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}

func (repo *PostRepo) DeletePost(id primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err := repo.collection.DeleteOne(c, filter)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting post from database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully."})
}

func (repo *PostRepo) GetPostById(id primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}
	var post models.Post

	if err := repo.collection.FindOne(c, filter).Decode(&post); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (repo *PostRepo) GetPosts(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := repo.collection.Find(c, bson.M{})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var posts []models.Post
	if err := cursor.All(c, &posts); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (repo *PostRepo) UpdatePost(id primitive.ObjectID, post models.Post, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": post}

	if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating post in database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post updated successfully."})
}

func (repo *PostRepo) UpVotePost(id primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id, "upVotedBy": userId}
	var post models.Post

	err := repo.collection.FindOne(c, filter).Decode(&post)
	if err == mongo.ErrNoDocuments {
		update := bson.M{"$inc": bson.M{"upVotes": 1},
			"$push": bson.M{"upVotedBy": userId},
		}
		filter = bson.M{"_id": id}
		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": "Post upvoted successfully."})
		}
	} else {
		update := bson.M{"$inc": bson.M{"upVotes": -1},
			"$pull": bson.M{"upVotedBy": userId},
		}
		filter = bson.M{"_id": id}
		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": "Upvote removed."})
		}
	}
}

func (repo *PostRepo) DownVote(id primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout((context.Background()), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id, "downVotedBy": userId}
	var post models.Post

	err := repo.collection.FindOne(c, filter).Decode(&post)
	if err == mongo.ErrNoDocuments {
		update := bson.M{"$inc": bson.M{"downVotes": 1},
			"$push": bson.M{"downVotedBy": userId},
		}

		filter = bson.M{"_id": id}

		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": "Post downvoted successfully."})
		}
	} else {
		update := bson.M{"$inc": bson.M{"downVotes": -1},
			"$pull": bson.M{"downVotedBy": userId},
		}

		filter = bson.M{"_id": id}

		if _, err := repo.collection.UpdateOne(c, filter, update); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"message": "Downvote removed."})
		}
	}
}
