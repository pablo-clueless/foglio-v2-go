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
			Username: payload.Username,
		}

		created, err := h.service.CreateUser(user)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
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
			if err.Error() == "user not found" {
				lib.NotFound(ctx, err.Error(), "")
				return
			}
			if err.Error() == "invalid password" || err.Error() == "user not verified" {
				lib.Unauthorized(ctx, err.Error())
				return
			}
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "User signed in successfully", user)
	}
}

func (h *AuthHandler) RequestVerification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")
		if email == "" {
			lib.BadRequest(ctx, "email is required", "400")
			return
		}

		user, err := h.service.RequestVerification(email)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		ctx.SetCookie("verification_email", user.Email, 1800, "/", "localhost", false, true)
		lib.Success(ctx, "", user)
	}
}

func (h *AuthHandler) Verification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.VerificationDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user, err := h.service.Verification(payload.Otp)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "User verified", user)
	}
}

func (h *AuthHandler) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.ChangePasswordDto
		id := ctx.Param("id")

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ChangePassword(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Password changed successfully", nil)
	}
}

func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.ForgotPasswordDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ForgotPassword(payload.Email)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "A reset mail has been sent to you", nil)

	}
}

func (h *AuthHandler) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.ResetPasswordDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		err := h.service.ResetPassword(payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Password reset successfully", nil)
	}
}

func (h *AuthHandler) GetOAuthURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.Param("provider")

		url, err := h.service.GetOAuthURL(provider)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "OAuth URL generated successfully", url)
	}
}

func (h *AuthHandler) HandleOAuthCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.Param("provider")
		var payload dto.OAuthCallbackDto

		if err := ctx.ShouldBindQuery(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		response, err := h.service.HandleOAuthCallback(provider, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error, "+err.Error())
			return
		}

		lib.Success(ctx, "", response)
	}
}
