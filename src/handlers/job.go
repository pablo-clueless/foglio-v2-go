package handlers

import (
	"foglio/v2/src/database"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	service *services.JobService
}

func NewJobHandler() *JobHandler {
	return &JobHandler{
		service: services.NewJobService(database.GetDatabase()),
	}
}

func (h *JobHandler) CreateJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *JobHandler) UpdateJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *JobHandler) DeleteJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *JobHandler) GetJobs() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *JobHandler) GetJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *JobHandler) ApplyToJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
