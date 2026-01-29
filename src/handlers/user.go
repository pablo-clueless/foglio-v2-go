package handlers

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/dto"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"
	"log"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: services.NewUserService(database.GetDatabase()),
	}
}

func (h *UserHandler) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query dto.UserPagination

		if err := ctx.ShouldBindQuery(&query); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		users, err := h.service.GetUsers(query)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Users fetched successfully", users)
	}
}

func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		user, err := h.service.GetUser(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User fetched successfully", user)
	}
}

func (h *UserHandler) GetMe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)

		if id == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		user, err := h.service.GetUser(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User fetched successfully", user)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateUserDto
		id := ctx.Param("id")

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		user, err := h.service.UpdateUser(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "User updated successfully", user)
	}
}

func (h *UserHandler) UpdateAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		file, header, err := ctx.Request.FormFile("avatar")
		if err != nil {
			lib.BadRequest(ctx, "avatar field is required", "400")
			return
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()

		url, err := lib.UploadSingle(header, "foglio-images")
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		user, err := h.service.UpdateAvatar(id, url)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "Avatar updated successfully", user)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteUser(id)
		if err != nil {
			lib.InternalServerError(ctx, "Internal server error,"+err.Error())
			return
		}

		lib.Success(ctx, "user deleted successfully", nil)
	}
}
