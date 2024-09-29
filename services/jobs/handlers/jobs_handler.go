package handlers

import (
	"net/http"

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

	id := primitive.NewObjectID()
	job.ID = id

	job.AppliedBy = []primitive.ObjectID{}
	job.LikedBy = []primitive.ObjectID{}

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
	var application data.JobApplication
	
	if err := ctx.ShouldBindJSON(&application); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := primitive.NewObjectID()
	application.Id = id

	if err := application.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jh.repo.ApplyForJob(ctx, application)
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

func (jh *JobsHandler) JobsLikedByMe(ctx *gin.Context) {
	userIdRaw := ctx.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.JobsLikedByMe(ctx, userId)
}

func (jh *JobsHandler) JobsAppliedByMe(ctx *gin.Context) {
	userIdRaw := ctx.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.JobsAppliedByMe(ctx, userId)
}

func (jh *JobsHandler) JobsByMe(ctx *gin.Context) {
	userIdRaw := ctx.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdRaw)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.JobsByMe(ctx, userId)
}

func (jh *JobsHandler) GetJobApplicationsByJob(ctx *gin.Context) {
	id := ctx.Param("id")
	jobId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.GetApplicationsForJob(ctx, jobId)
}

func (jh *JobsHandler) GetJobApplicationsByUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid id"})
		return
	}

	jh.repo.GetApplicationsByUser(ctx, userId)
}