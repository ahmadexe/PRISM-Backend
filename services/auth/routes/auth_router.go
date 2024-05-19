package routes

import (
	firebase "firebase.google.com/go"
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/ahmadexe/prism-backend/services/auth/middlewares"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authHandler *handlers.AuthHandler
	router      *gin.Engine
}

func InitAuthRouter(authHandler *handlers.AuthHandler, router *gin.Engine) *AuthRouter {
	return &AuthRouter{authHandler: authHandler, router: router}
}

func (router *AuthRouter) SetupRoutes(app *firebase.App) {
	auth := router.router.Group("/v1")
	auth.Use(middlewares.VerifyUser)
	{
		auth.POST("/users", router.authHandler.AddUser)
		auth.GET("/users/:id", router.authHandler.GetUserById)
		auth.GET("/users/primitive/:id", router.authHandler.GetById)
		auth.PUT("/users", router.authHandler.UpdateUser)
	}
}
