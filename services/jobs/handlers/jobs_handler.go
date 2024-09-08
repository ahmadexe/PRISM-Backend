package handlers

import (
	"github.com/ahmadexe/prism-backend/services/jobs/data"
	"github.com/ahmadexe/prism-backend/services/jobs/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobsHandler struct {
	repo *repository.JobsRepo
}

func InitJobHandler(repo *repository.JobsRepo) *JobsHandler {
	return &JobsHandler{repo: repo}
}

func (jh *JobsHandler) CreateJob(ctx *gin.Context) {
	var job data.Job
	if err := ctx.ShouldBindJSON(&job); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := job.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jh.repo.CreateJob(ctx, &job)
}

func (jh *JobsHandler) GetJobs(ctx *gin.Context) {
	jh.repo.GetJobs(ctx)
}

func (jh *JobsHandler) GetJob(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.GetJob(ctx, id)
}

func (jh *JobsHandler) ApplyJob(ctx *gin.Context) {
	var application data.Request
	if err := ctx.ShouldBindJSON(&application); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := application.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jh.repo.ApplyForJob(ctx, application.ID, application.UserId)
}

func (jh *JobsHandler) LikeJob(ctx *gin.Context) {
	var application data.Request
	if err := ctx.ShouldBindJSON(&application); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := application.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jh.repo.ToggleLikeOnJob(ctx, application.ID, application.UserId)
}

func (jh *JobsHandler) HireApplicant(ctx *gin.Context) {
	var application data.Request
	if err := ctx.ShouldBindJSON(&application); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := application.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jh.repo.HireForJob(ctx, application.ID, application.UserId)
}

func (jh *JobsHandler) UpdateJob(ctx *gin.Context) {
	var job data.Job
	if err := ctx.ShouldBindJSON(&job); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := job.Validate(); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jh.repo.UpdateJob(ctx, &job)
}

func (jh *JobsHandler) DeleteJob(ctx *gin.Context) {
	idRaw := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(idRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.DeleteJob(ctx, id)
}