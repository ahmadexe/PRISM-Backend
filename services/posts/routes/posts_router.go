package routes

import (
	"github.com/ahmadexe/prism-backend/services/posts/handlers"
	"github.com/gin-gonic/gin"
)

type PostsRouter struct {
	postsHandler *handlers.PostHandler
	router *gin.Engine
}

func InitPostsRouter(postsHandler *handlers.PostHandler, router *gin.Engine) *PostsRouter {
	return &PostsRouter{postsHandler: postsHandler, router: router}
}

func (r *PostsRouter) SetupRoutes() {
	posts := r.router.Group("/v1")
	{
		posts.POST("/posts", r.postsHandler.AddPost)
	}
}

