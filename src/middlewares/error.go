package middlewares

import (
	"fmt"
	"foglio/v2/src/lib"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			switch e := err.Err.(type) {
			case *lib.ApiError:
				ctx.AbortWithStatusJSON(e.Status, e)
			case validator.ValidationErrors:
				for _, fieldErr := range e {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"field":  fieldErr.Field(),
						"error":  fmt.Sprintf("Validation failed on tag '%s'", fieldErr.Tag()),
						"values": fieldErr.Param(),
					})
					return
				}
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, lib.NewApiErrror(e.Error(), http.StatusInternalServerError))
			}
		}
	}
}
