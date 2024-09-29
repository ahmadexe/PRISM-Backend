package routes

import (
	"github.com/ahmadexe/prism-backend/services/jobs/handlers"
	// "github.com/ahmadexe/prism-backend/services/jobs/middlewares"
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
	// jobs.Use(middlewares.VerifyUser)
	{
		jobs.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "Welcome to Prism Jobs Service"})
		})
		jobs.POST("/jobs", router.jobsHandler.CreateJob)
		jobs.GET("/jobs", router.jobsHandler.GetJobs)
		jobs.GET("/jobs/:id", router.jobsHandler.GetJob)
		jobs.POST("/jobs/apply", router.jobsHandler.ApplyJob)
		jobs.PUT("/jobs", router.jobsHandler.UpdateJob)
		jobs.DELETE("/jobs/:id", router.jobsHandler.DeleteJob)
		jobs.PUT("/jobs/like", router.jobsHandler.LikeJob)
		jobs.PUT("/jobs/hire", router.jobsHandler.HireApplicant)
		jobs.GET("/jobs/applied/:id", router.jobsHandler.JobsAppliedByMe)
		jobs.GET("/jobs/liked/:id", router.jobsHandler.JobsLikedByMe)
		jobs.GET("/jobs/user/:id", router.jobsHandler.JobsByMe)
	}
}
