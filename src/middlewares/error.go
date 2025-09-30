package middlewares

import (
	"foglio/v2/src/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	Message string                 `json:"message"`
	Errors  []ValidationFieldError `json:"errors"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last()

		switch e := err.Err.(type) {
		case *lib.ApiError:
			ctx.AbortWithStatusJSON(e.Status, e)

		case validator.ValidationErrors:
			handleValidationErrors(ctx, e)

		default:
			handleGenericError(ctx, e)
		}
	}
}

func handleValidationErrors(ctx *gin.Context, validationErrs validator.ValidationErrors) {
	fieldErrors := make([]ValidationFieldError, 0, len(validationErrs))

	for _, fieldErr := range validationErrs {
		fieldErrors = append(fieldErrors, ValidationFieldError{
			Field:   fieldErr.Field(),
			Message: getValidationMessage(fieldErr),
			Value:   fieldErr.Param(),
		})
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, ValidationErrorResponse{
		Message: "Validation failed",
		Errors:  fieldErrors,
	})
}

func getValidationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + fieldErr.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + fieldErr.Param() + ")"
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	case "gte":
		return "Value must be greater than or equal to " + fieldErr.Param()
	case "lte":
		return "Value must be less than or equal to " + fieldErr.Param()
	case "len":
		return "Length must be " + fieldErr.Param()
	case "oneof":
		return "Value must be one of: " + fieldErr.Param()
	default:
		return "Validation failed on '" + fieldErr.Tag() + "'"
	}
}

func handleGenericError(ctx *gin.Context, err error) {
	apiErr := lib.NewApiErrror(err.Error(), http.StatusInternalServerError)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, apiErr)
}
