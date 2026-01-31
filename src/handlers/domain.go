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

type DomainHandler struct {
	service *services.DomainService
}

func NewDomainHandler() *DomainHandler {
	return &DomainHandler{
		service: services.NewDomainService(database.GetDatabase()),
	}
}

func (h *DomainHandler) GetDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		domain, err := h.service.GetDomain(userId)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Domain configuration retrieved", domain)
	}
}

func (h *DomainHandler) CheckSubdomainAvailability() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subdomain := ctx.Param("subdomain")
		if subdomain == "" {
			lib.BadRequest(ctx, "subdomain is required", "")
			return
		}

		result, err := h.service.CheckSubdomainAvailability(subdomain)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Subdomain availability checked", result)
	}
}

func (h *DomainHandler) ClaimSubdomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.ClaimSubdomainDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		domain, err := h.service.ClaimSubdomain(userId, payload)
		if err != nil {
			handleDomainError(ctx, err)
			return
		}

		lib.Created(ctx, "Subdomain claimed successfully", domain)
	}
}

func (h *DomainHandler) UpdateSubdomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.ClaimSubdomainDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		domain, err := h.service.UpdateSubdomain(userId, payload)
		if err != nil {
			handleDomainError(ctx, err)
			return
		}

		lib.Success(ctx, "Subdomain updated successfully", domain)
	}
}

func (h *DomainHandler) SetCustomDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var payload dto.SetCustomDomainDto
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "")
			return
		}

		domain, err := h.service.SetCustomDomain(userId, payload)
		if err != nil {
			handleDomainError(ctx, err)
			return
		}

		lib.Created(ctx, "Custom domain configured. Please add the DNS records to verify ownership.", domain)
	}
}

func (h *DomainHandler) VerifyCustomDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		domain, err := h.service.VerifyCustomDomain(userId)
		if err != nil {
			handleDomainError(ctx, err)
			return
		}

		message := "DNS verification in progress"
		if domain.CustomDomainStatus == "VERIFIED" {
			message = "Custom domain verified successfully"
		}

		lib.Success(ctx, message, domain)
	}
}

func (h *DomainHandler) RemoveCustomDomain() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetString(config.AppConfig.CurrentUserId)
		if userId == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		domain, err := h.service.RemoveCustomDomain(userId)
		if err != nil {
			handleDomainError(ctx, err)
			return
		}

		lib.Success(ctx, "Custom domain removed successfully", domain)
	}
}

func handleDomainError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrSubdomainTaken):
		lib.BadRequest(ctx, "This subdomain is already taken", "SUBDOMAIN_TAKEN")
	case errors.Is(err, services.ErrSubdomainInvalid):
		lib.BadRequest(ctx, "Subdomain contains invalid characters. Use only lowercase letters, numbers, and hyphens.", "SUBDOMAIN_INVALID")
	case errors.Is(err, services.ErrSubdomainReserved):
		lib.BadRequest(ctx, "This subdomain is reserved and cannot be used", "SUBDOMAIN_RESERVED")
	case errors.Is(err, services.ErrCustomDomainRequired):
		lib.Forbidden(ctx, "Custom domains require a paid subscription. Please upgrade to use this feature.")
	case errors.Is(err, services.ErrDomainAlreadySet):
		lib.BadRequest(ctx, "You already have a subdomain. Use the update endpoint to change it.", "DOMAIN_ALREADY_SET")
	case errors.Is(err, services.ErrNoCustomDomain):
		lib.BadRequest(ctx, "No custom domain configured", "NO_CUSTOM_DOMAIN")
	default:
		lib.InternalServerError(ctx, err.Error())
	}
}
