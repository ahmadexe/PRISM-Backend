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
		posts.DELETE("/posts/:id", r.postsHandler.DeletePost)
		posts.GET("/posts/:id", r.postsHandler.GetPostById)
		posts.GET("/posts", r.postsHandler.GetPosts)
		posts.PUT("/posts/:id", r.postsHandler.UpdatePost)
	}
}

