package routes

import (
	"github.com/ahmadexe/prism-backend/services/jobs/handlers"
	"github.com/gin-gonic/gin"
)

type JobsRouter struct {
	jobsHandler *handlers.JobsHandler
	router      *gin.Engine
}

func InitJobsRouter(jobsHandler *handlers.JobsHandler, router *gin.Engine) *JobsRouter {
	return &JobsRouter{jobsHandler: jobsHandler, router: router}
}

func (router *JobsRouter) SetupRoutes() {
	jobs := router.router.Group("/v1")
	{
		jobs.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "Welcome to Prism Jobs Service"})
		})
	}
}
