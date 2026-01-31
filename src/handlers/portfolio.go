package handlers

import (
	"errors"
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	service *services.PortfolioService
}

func NewPortfolioHandler() *PortfolioHandler {
	return &PortfolioHandler{
		service: services.NewPortfolioService(database.GetDatabase()),
	}
}

// CreatePortfolio creates a new portfolio for the authenticated user
func (h *PortfolioHandler) CreatePortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.CreatePortfolioDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		portfolio, err := h.service.CreatePortfolio(userId, payload)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Created(ctx, "Portfolio created successfully", h.service.ToResponse(portfolio))
	}
}

// GetPortfolio returns the authenticated user's portfolio
func (h *PortfolioHandler) GetPortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		portfolio, err := h.service.GetPortfolio(userId)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio retrieved successfully", h.service.ToResponse(portfolio))
	}
}

// GetPortfolioBySlug returns a public portfolio by slug
func (h *PortfolioHandler) GetPortfolioBySlug() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slug := ctx.Param("slug")
		if slug == "" {
			lib.BadRequest(ctx, "slug is required", "")
			return
		}

		portfolio, user, err := h.service.GetPortfolioBySlug(slug)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio retrieved successfully", h.service.ToPublicResponse(portfolio, user))
	}
}

// UpdatePortfolio updates the authenticated user's portfolio
func (h *PortfolioHandler) UpdatePortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.UpdatePortfolioDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		portfolio, err := h.service.UpdatePortfolio(userId, payload)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio updated successfully", h.service.ToResponse(portfolio))
	}
}

// DeletePortfolio deletes the authenticated user's portfolio
func (h *PortfolioHandler) DeletePortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		if err := h.service.DeletePortfolio(userId); err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio deleted successfully", nil)
	}
}

// PublishPortfolio publishes the portfolio
func (h *PortfolioHandler) PublishPortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		portfolio, err := h.service.PublishPortfolio(userId)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio published successfully", h.service.ToResponse(portfolio))
	}
}

// UnpublishPortfolio unpublishes the portfolio
func (h *PortfolioHandler) UnpublishPortfolio() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		portfolio, err := h.service.UnpublishPortfolio(userId)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Portfolio unpublished successfully", h.service.ToResponse(portfolio))
	}
}

// CreateSection creates a new section in the portfolio
func (h *PortfolioHandler) CreateSection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.CreatePortfolioSectionDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		section, err := h.service.CreateSection(userId, payload)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Created(ctx, "Section created successfully", section)
	}
}

// UpdateSection updates a section
func (h *PortfolioHandler) UpdateSection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		sectionId := ctx.Param("sectionId")
		if sectionId == "" {
			lib.BadRequest(ctx, "section ID is required", "")
			return
		}

		var payload dto.UpdatePortfolioSectionDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		section, err := h.service.UpdateSection(userId, sectionId, payload)
		if err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Section updated successfully", section)
	}
}

// DeleteSection deletes a section
func (h *PortfolioHandler) DeleteSection() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		sectionId := ctx.Param("sectionId")
		if sectionId == "" {
			lib.BadRequest(ctx, "section ID is required", "")
			return
		}

		if err := h.service.DeleteSection(userId, sectionId); err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Section deleted successfully", nil)
	}
}

// ReorderSections reorders sections
func (h *PortfolioHandler) ReorderSections() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.ReorderSectionsDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		if err := h.service.ReorderSections(userId, payload); err != nil {
			handlePortfolioError(ctx, err)
			return
		}

		lib.Success(ctx, "Sections reordered successfully", nil)
	}
}

func handlePortfolioError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrPortfolioNotFound):
		lib.NotFound(ctx, "Portfolio not found", "PORTFOLIO_NOT_FOUND")
	case errors.Is(err, services.ErrPortfolioExists):
		lib.BadRequest(ctx, "You already have a portfolio", "PORTFOLIO_EXISTS")
	case errors.Is(err, services.ErrSlugTaken):
		lib.BadRequest(ctx, "This slug is already taken", "SLUG_TAKEN")
	case errors.Is(err, services.ErrSlugInvalid):
		lib.BadRequest(ctx, "Slug must be 3-50 characters, lowercase letters, numbers, and hyphens only", "SLUG_INVALID")
	case errors.Is(err, services.ErrSectionNotFound):
		lib.NotFound(ctx, "Section not found", "SECTION_NOT_FOUND")
	case errors.Is(err, services.ErrUnauthorized):
		lib.Forbidden(ctx, "You are not authorized to perform this action")
	default:
		lib.InternalServerError(ctx, err.Error())
	}
}
