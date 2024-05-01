package routes

import (
	firebase "firebase.google.com/go"
	"github.com/ahmadexe/prism-backend/middlewares"
	"github.com/ahmadexe/prism-backend/services/auth/handlers"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authHandler *handlers.AuthHandler
	router *gin.Engine
}

func InitAuthRouter(authHandler *handlers.AuthHandler, router *gin.Engine) *AuthRouter {
	return &AuthRouter{authHandler: authHandler, router: router}
}

func (router *AuthRouter) SetupRoutes(app *firebase.App) {
	authMiddleware := middlewares.InitAuthMiddleware(app)
	auth := router.router.Group("/auth")
	auth.Use(authMiddleware.VerifyUser)
	{
		auth.POST("/users", router.authHandler.AddUser)
		auth.GET("/users/:id", router.authHandler.GetUserById)
		auth.PUT("/users", router.authHandler.UpdateUser)
	}
}
