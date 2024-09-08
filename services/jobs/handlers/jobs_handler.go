package handlers

import "github.com/ahmadexe/prism-backend/services/jobs/repository"

type JobsHandler struct {
	repo *repository.JobsRepo
}

func InitJobHandler(repo *repository.JobsRepo) *JobsHandler {
	return &JobsHandler{repo: repo}
}