package lib

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      string    `json:"code"`
	Path      string    `json:"path,omitempty"`
	Method    string    `json:"method,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *string     `json:"error,omitempty"`
	Code      string      `json:"code,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

const (
	NotFoundCode         = "RESOURCE_NOT_FOUND"
	UserNotFoundCode     = "USER_NOT_FOUND"
	ProductNotFoundCode  = "PRODUCT_NOT_FOUND"
	OrderNotFoundCode    = "ORDER_NOT_FOUND"
	CategoryNotFoundCode = "CATEGORY_NOT_FOUND"
	ValidationErrorCode  = "VALIDATION_ERROR"
	InternalErrorCode    = "INTERNAL_ERROR"
	UnauthorizedCode     = "UNAUTHORIZED"
	ForbiddenCode        = "FORBIDDEN"
)

func GlobalNotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		NotFound(ctx, "The requested endpoint does not exist", "")
	}
}

func NotFound(ctx *gin.Context, message string, code string) {
	if message == "" {
		message = "Resource not found"
	}
	if code == "" {
		code = NotFoundCode
	}

	response := ErrorResponse{
		Success:   false,
		Error:     "Not Found",
		Message:   message,
		Code:      code,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusNotFound, response)
	ctx.Abort()
}

func UserNotFound(ctx *gin.Context) {
	NotFound(ctx, "User not found", UserNotFoundCode)
}

func ProductNotFound(ctx *gin.Context) {
	NotFound(ctx, "Product not found", ProductNotFoundCode)
}

func OrderNotFound(ctx *gin.Context) {
	NotFound(ctx, "Order not found", OrderNotFoundCode)
}

func CategoryNotFound(ctx *gin.Context) {
	NotFound(ctx, "Category not found", CategoryNotFoundCode)
}

func Success(ctx *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusOK, response)
}

func Created(ctx *gin.Context, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusCreated, response)
}

func BadRequest(ctx *gin.Context, message string, code string) {
	if code == "" {
		code = ValidationErrorCode
	}

	errorMsg := "Bad Request"
	response := ErrorResponse{
		Success:   false,
		Error:     errorMsg,
		Message:   message,
		Code:      code,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusBadRequest, response)
	ctx.Abort()
}

func InternalServerError(ctx *gin.Context, message string) {
	if message == "" {
		message = "Internal server error occurred"
	}

	response := ErrorResponse{
		Success:   false,
		Error:     "Internal Server Error",
		Message:   message,
		Code:      InternalErrorCode,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusInternalServerError, response)
	ctx.Abort()
}

func Unauthorized(ctx *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}

	response := ErrorResponse{
		Success:   false,
		Error:     "Unauthorized",
		Message:   message,
		Code:      UnauthorizedCode,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusUnauthorized, response)
	ctx.Abort()
}

func Forbidden(ctx *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}

	response := ErrorResponse{
		Success:   false,
		Error:     "Forbidden",
		Message:   message,
		Code:      ForbiddenCode,
		Path:      ctx.Request.URL.Path,
		Method:    ctx.Request.Method,
		Timestamp: time.Now().UTC(),
	}

	ctx.JSON(http.StatusForbidden, response)
	ctx.Abort()
}

func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		if recovered != nil {
			InternalServerError(ctx, "An unexpected error occurred")
		}
	})
}
