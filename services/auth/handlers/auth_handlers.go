package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/ahmadexe/prism-backend/services/auth/data"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepo
}

func InitAuthHandler(authRepo *repositories.AuthRepo) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

func (handler *AuthHandler) AddUser(ctx *gin.Context) {
	var user data.AuthData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	if err := user.Validate(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid data."})
		return
	}

	user.CreatedAt = time.Now().UnixMicro()
	user.Id = primitive.NewObjectID()
	user.Followers = []primitive.ObjectID{}
	user.Following = []primitive.ObjectID{}

	handler.authRepo.AddUser(user, ctx)
}

func (handler *AuthHandler) GetUserByUid(ctx *gin.Context) {
	handler.authRepo.GetUserById(ctx)
}

func (handler *AuthHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	handler.authRepo.GetById(objectId, ctx)
}

func (handler *AuthHandler) UpdateUser(ctx *gin.Context) {
	var user data.AuthData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Validate(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.authRepo.UpdateUser(user, ctx)
}

func (handler *AuthHandler) ToggleFollowRequest(ctx *gin.Context) {
	var followReq data.FollowRequest

	if err := ctx.ShouldBindJSON(&followReq); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	if err := followReq.Validate(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.authRepo.ToggleFollow(followReq, ctx)
}


func (handler *AuthHandler) ToggleIsServiceProvider(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	handler.authRepo.ToggleIsServiceProvider(objectId, ctx)
}

func (handler *AuthHandler) UpdateDeviceToken(ctx *gin.Context) {
	var tokenReq data.TokenRequest

	if err := ctx.ShouldBindJSON(&tokenReq); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	if err := tokenReq.Validate(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.authRepo.UpdateDeviceToken(ctx, tokenReq)
}

func (handler *AuthHandler) ToggleIsSupercharged(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	handler.authRepo.ToggleSupercharged(objectId, ctx)
}

func (handler *AuthHandler) GetFollowers(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	handler.authRepo.GetFollowers(objectId, ctx)
}

func (handler *AuthHandler) GetFollowing(ctx *gin.Context) {
	id := ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	handler.authRepo.GetFollowing(objectId, ctx)
}