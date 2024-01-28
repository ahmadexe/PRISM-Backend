package routes

import (
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	AuthHandler *handlers.AuthHandler
	Router *gin.Engine
}

func InitAuthRouter(authHandler *handlers.AuthHandler, router *gin.Engine) *AuthRouter {
	return &AuthRouter{AuthHandler: authHandler, Router: router}
}

func (router *AuthRouter) SetupRoutes() {
	auth := router.Router.Group("/auth")
	{
		auth.POST("/add-user", router.AuthHandler.AddUser)
	}
}
