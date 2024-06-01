package handlers

import (
	"net/http"

	"github.com/ahmadexe/prism-backend/services/posts/data"
	"github.com/ahmadexe/prism-backend/services/posts/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	repo *repositories.PostRepo
}

func InitCommentHandler(repo *repositories.PostRepo) *CommentHandler {
	return &CommentHandler{repo: repo}
}

func (h *CommentHandler) AddComment(ctx *gin.Context) {
	var comment data.Comment

	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	if err := comment.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	comment.Id = primitive.NewObjectID()
	h.repo.AddComment(comment, comment.PostId, comment.UserId, ctx)
}

func (h *CommentHandler) DeleteComment(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid comment id."})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid comment id."})
		return
	}

	postId := ctx.Param("postId")
	if postId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	postObjectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	userId := ctx.Param("userId")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid user id."})
		return
	}
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid user id."})
		return
	}

	h.repo.DeleteComment(objectId, postObjectId, userObjectId, ctx)
}

func (h *CommentHandler) GetComments(ctx *gin.Context) {
	postId := ctx.Param("id")

	if postId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a valid post id."})
		return
	}

	h.repo.GetComments(objectId, ctx)
}
