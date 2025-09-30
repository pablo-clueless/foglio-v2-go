package handlers

import (
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
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
	return func(ctx *gin.Context) {
		var payload dto.CreateJobDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		job, err := h.service.CreateJob(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Job created successfully", job)
	}
}

func (h *JobHandler) UpdateJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateJobDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		job, err := h.service.UpdateJob(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Job updated successfully", job)
	}
}

func (h *JobHandler) DeleteJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteJob(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Job updated successfully", nil)
	}
}

func (h *JobHandler) GetJobs() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query dto.JobPagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		jobs, err := h.service.GetJobs(query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Users fetched successfully", jobs)
	}
}

func (h *JobHandler) GetJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		job, err := h.service.GetJob(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Job updated successfully", job)
	}
}

func (h *JobHandler) ApplyToJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		jobId := ctx.Param("jobId")
		var payload dto.JobApplicationDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ApplyToJob(userId, jobId, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Job application success", nil)
	}
}
