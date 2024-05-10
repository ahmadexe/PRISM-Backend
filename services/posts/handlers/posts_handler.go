package handlers

import (
	"net/http"

	"github.com/ahmadexe/prism-backend/services/posts/models"
	"github.com/ahmadexe/prism-backend/services/posts/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHandler struct {
	repo *repositories.PostRepo
}

func InitPostHandler(repo *repositories.PostRepo) *PostHandler {
	return &PostHandler{repo: repo}
}

func (handler *PostHandler) AddPost(ctx *gin.Context) {
	var post models.Post

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	if err := post.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	post.Id = primitive.NewObjectID()
	post.UpVotedBy = []string{}
	post.DownVotedBy = []string{}
	post.CommentedBy = []string{}

	handler.repo.AddPost(post, ctx)
}

func (handler *PostHandler) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	handler.repo.DeletePost(objectId, ctx)
}

func (handler *PostHandler) GetPostById(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	handler.repo.GetPostById(objectId, ctx)
}

func (handler *PostHandler) GetPosts(ctx *gin.Context) {
	handler.repo.GetPosts(ctx)
}

func (handler *PostHandler) UpdatePost(ctx *gin.Context) {
	var post models.Post

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	if err := post.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	handler.repo.UpdatePost(objectId, post, ctx)
}
