package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/models"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	service *services.JobService
}

func NewJobHandler() *JobHandler {
	return &JobHandler{
		service: services.NewJobService(database.GetDatabase(), services.NewNotificationService(database.GetDatabase(), lib.NewHub())),
	}
}

func (h *JobHandler) CreateJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateJobDto
		userId := ctx.GetString(config.AppConfig.CurrentUserId)

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		job, err := h.service.CreateJob(userId, payload)
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

		lib.Success(ctx, "Jobs fetched successfully", jobs)
	}
}

func (h *JobHandler) GetJobsByUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		var query dto.Pagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		jobs, err := h.service.GetJobsByUser(id, query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Jobs fetched successfully", jobs)
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
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		jobId := ctx.Param("id")
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

		lib.Success(ctx, "Job application successful", nil)
	}
}

func (h *JobHandler) GetApplicationsByUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		var query dto.Pagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		applications, err := h.service.GetApplicationsByUser(id, query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Applications fetched successfully", applications)

	}
}

func (h *JobHandler) GetApplicationsByJob() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		jobId := ctx.Param("id")
		var query dto.JobApplicationPagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		applications, err := h.service.GetApplicationsByJob(id, jobId, query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Applications fetched successfully", applications)
	}
}

func (h *JobHandler) GetApplicationsByRecruiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		var query dto.JobApplicationPagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		applications, err := h.service.GetApplicationsByRecruiter(id, query)
		if err != nil {
			if err.Error() == "only recruiters can view applications" {
				lib.Forbidden(ctx, err.Error())
				return
			}
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Applications fetched successfully", applications)
	}
}

func (h *JobHandler) AcceptApplication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		applicationId := ctx.Param("id")
		var payload dto.ApplicationStatusDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		application, err := h.service.AcceptApplication(id, applicationId, payload.Reason)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Application accepted successfully", application)
	}
}

func (h *JobHandler) RejectApplication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		applicationId := ctx.Param("id")
		var payload dto.ApplicationStatusDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		application, err := h.service.RejectApplication(id, applicationId, payload.Reason)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Application rejected successfully", application)
	}
}

func (h *JobHandler) ReviewApplication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		applicationId := ctx.Param("id")
		var payload dto.ApplicationStatusDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		application, err := h.service.ReviewApplication(id, applicationId, payload.Reason)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Application marked as reviewed", application)
	}
}

func (h *JobHandler) HireApplication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		applicationId := ctx.Param("id")
		var payload dto.ApplicationStatusDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		application, err := h.service.HireApplication(id, applicationId, payload.Reason)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Applicant hired successfully", application)
	}
}

func (h *JobHandler) GetApplication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		applicationId := ctx.Param("applicationId")

		application, err := h.service.GetApplicationById(applicationId)
		if err != nil {
			if err.Error() == "application not found" {
				lib.NotFound(ctx, err.Error(), "404")
				return
			}
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Application fetched successfully", application)
	}
}

func (h *JobHandler) AddComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		jobId := ctx.Param("id")
		var payload dto.CommentDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		created, err := h.service.AddComment(id, jobId, payload.Content)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Comment added successfully", created)
	}
}

func (h *JobHandler) DeleteComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		commentId := ctx.Param("id")

		err := h.service.DeleteComment(commentId, id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Comment deleted successfully", nil)
	}
}

func (h *JobHandler) AddReaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		jobId := ctx.Param("id")
		reactionType := ctx.Param("reaction")

		reaction, err := h.service.AddReaction(id, jobId, models.ReactionType(reactionType))
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Reaction added successfully", reaction)
	}
}

func (h *JobHandler) RemoveReaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		jobId := ctx.Param("id")

		err := h.service.RemoveReaction(id, jobId)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Reaction removed successfully", nil)
	}
}
