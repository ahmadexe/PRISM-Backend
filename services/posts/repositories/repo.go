package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ahmadexe/prism-backend/services/posts/data"
	"github.com/ahmadexe/prism-backend/services/posts/network"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepo struct {
	postCollection    *mongo.Collection
	commentCollection *mongo.Collection
}

func InitPostRepo(client *mongo.Client) *PostRepo {
	postCollection := client.Database("posts-db").Collection("posts")
	commentCollection := client.Database("posts-db").Collection("comments")
	return &PostRepo{postCollection: postCollection, commentCollection: commentCollection}
}

func (repo *PostRepo) AddPost(post data.Post, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := repo.postCollection.InsertOne(c, post)

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

	_, err := repo.postCollection.DeleteOne(c, filter)
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
	var post data.Post

	if err := repo.postCollection.FindOne(c, filter).Decode(&post); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, post)
}

func (repo *PostRepo) GetPosts(ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := repo.postCollection.Find(c, bson.M{})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var posts []data.Post
	if err := cursor.All(c, &posts); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (repo *PostRepo) UpdatePost(id primitive.ObjectID, post data.Post, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": post}

	var updatedPost data.Post

	if err := repo.postCollection.FindOneAndUpdate(c, filter, update).Decode(&updatedPost); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating post in database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post updated successfully.", "data": updatedPost})
}

func (repo *PostRepo) UpVotePost(id primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id, "upVotedBy": userId}
	var post data.Post

	err := repo.postCollection.FindOne(c, filter).Decode(&post)
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	if err == mongo.ErrNoDocuments {
		// Either there are no upvotes and no downvotes or the user has downvoted the post

		filter = bson.M{"_id": id, "downVotedBy": userId}
		err = repo.postCollection.FindOne(c, filter).Decode(&post)
		var update bson.M
		if err == mongo.ErrNoDocuments {
			// User has not voted on the post
			update = bson.M{
				"$inc":  bson.M{"upVotes": 1},
				"$push": bson.M{"upVotedBy": userId},
			}
		} else {
			// User has downvoted the post
			update = bson.M{
				"$inc":  bson.M{"upVotes": 1, "downVotes": -1},
				"$push": bson.M{"upVotedBy": userId},
				"$pull": bson.M{"downVotedBy": userId},
			}
		}

		filter = bson.M{"_id": id}
		err = repo.postCollection.FindOneAndUpdate(c, filter, update, options).Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Post not upvoted. Please try again later."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": post, "message": "Post upvoted successfully."})
	} else {
		// User has upvoted the post
		update := bson.M{"$inc": bson.M{"upVotes": -1},
			"$pull": bson.M{"upVotedBy": userId},
		}

		filter = bson.M{"_id": id}
		err = repo.postCollection.FindOneAndUpdate(c, filter, update, options).Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Upvote not removed. Please try again later."})
			return
		}
		ctx.JSON(200, gin.H{"data": post, "message": "Upvote removed."})
	}
}

func (repo *PostRepo) DownVote(id primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout((context.Background()), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id, "downVotedBy": userId}
	var post data.Post

	err := repo.postCollection.FindOne(c, filter).Decode(&post)

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err == mongo.ErrNoDocuments {
		// Either there are no upvotes and no downvotes or the user has upvoted the post

		filter = bson.M{"_id": id, "upVotedBy": userId}
		var update bson.M
		err = repo.postCollection.FindOne(c, filter).Decode(&post)
		if err == mongo.ErrNoDocuments {
			// User has not voted on the post
			update = bson.M{
				"$inc":  bson.M{"downVotes": 1},
				"$push": bson.M{"downVotedBy": userId},
			}
		} else {
			// User has upvoted the post
			update = bson.M{
				"$inc":  bson.M{"downVotes": 1, "upVotes": -1},
				"$push": bson.M{"downVotedBy": userId},
				"$pull": bson.M{"upVotedBy": userId},
			}
		}

		filter = bson.M{"_id": id}
		err = repo.postCollection.FindOneAndUpdate(c, filter, update, options).Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Post not downvoted. Please try again later."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": post, "message": "Post downvoted successfully."})
	} else {
		// User has downvoted the post
		update := bson.M{"$inc": bson.M{"downVotes": -1},
			"$pull": bson.M{"downVotedBy": userId},
		}

		filter = bson.M{"_id": id}
		err = repo.postCollection.FindOneAndUpdate(c, filter, update, options).Decode(&post)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Downvote not removed. Please try again later."})
			return
		}
		ctx.JSON(200, gin.H{"data": post, "message": "Downvote removed."})
	}
}

func (repo *PostRepo) AddComment(comment data.Comment, postId primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := repo.commentCollection.InsertOne(c, comment)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding comment to database. Please try again later."})
		return
	}

	filter := bson.M{"_id": postId}
	update := bson.M{"$inc": bson.M{"noOfComments": 1}, "$push": bson.M{"commentedBy": userId}}

	if _, err := repo.postCollection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding comment to database. Please try again later."})
		cid := comment.Id
		filter = bson.M{"_id": cid}
		_, err := repo.commentCollection.DeleteOne(c, filter)
		if err != nil {
			log.Println(err)
		}
		return
	}

	fmt.Println("Comment added successfully")

	ctx.JSON(http.StatusCreated, comment)
}

func (repo *PostRepo) GetComments(postId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cursor, err := repo.commentCollection.Find(c, bson.M{"postId": postId})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	var comments []data.Comment
	if err := cursor.All(c, &comments); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (repo *PostRepo) DeleteComment(id primitive.ObjectID, postId primitive.ObjectID, userId primitive.ObjectID, ctx *gin.Context) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err := repo.commentCollection.DeleteOne(c, filter)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting comment from database. Please try again later."})
		return
	}

	filter = bson.M{"_id": postId}
	update := bson.M{"$inc": bson.M{"noOfComments": -1}, "$pull": bson.M{"commentedBy": userId}}

	if _, err := repo.postCollection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting comment from database. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully."})
}

func (repo *PostRepo) Report(ctx *gin.Context, req data.ReportRequest) {
	c, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"_id": req.PostId}
	update := bson.M{"$inc": bson.M{"reportsRecord." + req.Type: 1, "totalReports": 1}}

	if _, err := repo.postCollection.UpdateOne(c, filter, update); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reporting post. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post reported successfully."})

	go func() {
		c, cancel := context.WithTimeout(context.Background(), time.Minute * 5)
		defer cancel()
		var post data.Post
		if err := repo.postCollection.FindOne(c, filter).Decode(&post); err != nil {
			log.Println(err)
		}
		if post.TotalReports >= 4 {
			max := -1
			highestReports := ""
			for key, value := range post.ReportsRecord {
				if value >= max {
					max = value
					highestReports = key
				}
			}
			
			req := map[string]interface{}{
				"url": post.ImageUrl,
				"post_id": 1,
			}

			reqRaw, err := json.Marshal(req) 
			if (err != nil) {
				log.Print(err)
			} 

			res := network.PostReq("https://prism-ml-5lyi.onrender.com/mod/image/"+highestReports, reqRaw)
			if (res["moderate"] == true) {
				update = bson.M{"$set": bson.M{"isBanned": true}}
				if _, err := repo.postCollection.UpdateOne(c, filter, update); err != nil {
					log.Println(err)
				}
			} 
		}
	}()
}

