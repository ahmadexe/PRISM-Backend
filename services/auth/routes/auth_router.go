package routes

import (
	"github.com/ahmadexe/prism-backend/middlewares"
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
	authMiddleware := middlewares.InitAuthMiddleware(router.AuthHandler.App)
	auth := router.Router.Group("/auth")
	auth.Use(authMiddleware.VerifyUser)
	{
		auth.POST("/users", router.AuthHandler.AddUser)
		auth.GET("/users/:id", router.AuthHandler.GetUserById)
		auth.PUT("/users/", router.AuthHandler.UpdateUser)
	}
}
