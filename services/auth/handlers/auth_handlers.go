package handlers

import (
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
	handler.AuthRepo.AddUser(ctx)
}
