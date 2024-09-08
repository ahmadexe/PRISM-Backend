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
		jobs.POST("/jobs", router.jobsHandler.CreateJob)
		jobs.GET("/jobs", router.jobsHandler.GetJobs)
		jobs.GET("/jobs/:id", router.jobsHandler.GetJob)
		jobs.PUT("/jobs/apply", router.jobsHandler.ApplyJob)
		jobs.PUT("/jobs", router.jobsHandler.UpdateJob)
		jobs.DELETE("/jobs/:id", router.jobsHandler.DeleteJob)
		jobs.PUT("/jobs/like", router.jobsHandler.LikeJob)
		jobs.PUT("/jobs/hire", router.jobsHandler.HireApplicant)
	}
}
