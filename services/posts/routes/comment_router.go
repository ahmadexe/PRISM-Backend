package routes

import (
	"github.com/ahmadexe/prism-backend/services/posts/handlers"
	"github.com/ahmadexe/prism-backend/services/posts/middlewares"
	"github.com/gin-gonic/gin"
)

type CommentRouter struct {
	commentHandler *handlers.CommentHandler
	router 	   *gin.Engine
}

func InitCommentRouter(commentHandler *handlers.CommentHandler, r *gin.Engine) *CommentRouter {
	return &CommentRouter{commentHandler: commentHandler, router: r}
}

func (r *CommentRouter) SetupRoutes() {
	comments := r.router.Group("/v1")
	comments.Use(middlewares.VerifyUser)
	{
		comments.POST("/comments", r.commentHandler.AddComment)
		comments.DELETE("/comments/:id/:postId/:userId", r.commentHandler.DeleteComment)
		comments.GET("/comments/:id", r.commentHandler.GetComments)
	}
}