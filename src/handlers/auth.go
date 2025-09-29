package handlers

import (
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(database.GetDatabase()),
	}
}

func (h *AuthHandler) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateUserDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user := dto.CreateUserDto{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: payload.Password,
		}

		created, err := h.service.CreateUser(user)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User created successfully", created)
	}
}

func (h *AuthHandler) Signin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.SigninDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user, err := h.service.Signin(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User signed in successfully", user)
	}
}

func (h *AuthHandler) Verification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var otp string

		if err := ctx.ShouldBind(&otp); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.Verification(otp)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User verified", nil)
	}
}

func (h *AuthHandler) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.ChangePasswordDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ChangePassword(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Password changed successfully", nil)
	}
}

func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var email string

		if err := ctx.ShouldBind(email); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ForgotPassword(email)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "A reset mail has been sent to you", nil)

	}
}

func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload dto.ResetPasswordDto

		if err := ctx.ShouldBind(payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ResetPassword(payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Password reset successfully", nil)
	}
}
