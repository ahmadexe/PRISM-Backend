package handlers

import (
	"net/http"
	"time"

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
	post.CreatedAt = time.Now().UnixMicro()

	handler.repo.AddPost(post, ctx)
}