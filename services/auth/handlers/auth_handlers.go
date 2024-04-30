package handlers

import (
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepo
	App *firebase.App
}

func InitAuthHandler(authRepo *repositories.AuthRepo, app *firebase.App) *AuthHandler {
	return &AuthHandler{authRepo: authRepo, App: app}
}

func (handler *AuthHandler) AddUser(ctx *gin.Context) {
	var user models.AuthData
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
	user.Followers = []string{}
	user.Following = []string{}

	handler.authRepo.AddUser(user, ctx)
}

func (handler *AuthHandler) GetUserById(ctx *gin.Context) {
	handler.authRepo.GetUserById(ctx)
}

func (handler *AuthHandler) UpdateUser(ctx *gin.Context) {
	var user models.AuthData
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