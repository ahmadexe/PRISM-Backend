package handlers

import (
	"github.com/ahmadexe/prism-backend/services/auth/models"
	"github.com/ahmadexe/prism-backend/services/auth/repositories"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthRepo *repositories.AuthRepo
}

func InitAuthHandler(authRepo *repositories.AuthRepo) *AuthHandler {
	return &AuthHandler{AuthRepo: authRepo}
}

func (handler *AuthHandler) AddUser(ctx *gin.Context) {
	var user models.AuthData
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := user.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	handler.AuthRepo.AddUser(user, ctx)
}

func (handler *AuthHandler) GetUserById(ctx *gin.Context) {
	handler.AuthRepo.GetUserById(ctx)
}